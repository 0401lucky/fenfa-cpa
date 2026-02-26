package proxy

import (
	"bytes"
	"cpa-distribution/common"
	"cpa-distribution/common/utils"
	"cpa-distribution/model"
	"cpa-distribution/service"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func ProxyHandler(c *gin.Context) {
	upstreamURL := common.CPAUpstreamURL
	upstreamKey := common.CPAUpstreamKey

	// Allow override from system settings
	if settingURL := model.GetSetting("cpa_upstream_url"); settingURL != "" {
		upstreamURL = settingURL
	}
	if settingKey := model.GetSetting("cpa_upstream_key"); settingKey != "" {
		upstreamKey = settingKey
	}

	if upstreamURL == "" || upstreamKey == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": gin.H{
				"message": "Upstream not configured",
				"type":    "server_error",
			},
		})
		return
	}

	target, err := url.Parse(upstreamURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Invalid upstream URL",
				"type":    "server_error",
			},
		})
		return
	}

	startTime := time.Now()

	// Read request body to extract model name
	var bodyBytes []byte
	var requestModel string
	var isStream bool
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var reqBody map[string]interface{}
		if json.Unmarshal(bodyBytes, &reqBody) == nil {
			if m, ok := reqBody["model"].(string); ok {
				requestModel = m
			}
			if s, ok := reqBody["stream"].(bool); ok {
				isStream = s
			}
		}
	}

	// Check allowed models
	if allowedModels, exists := c.Get("allowed_models"); exists {
		allowed := allowedModels.(string)
		if allowed != "" && requestModel != "" {
			modelAllowed := false
			for _, m := range strings.Split(allowed, ",") {
				if strings.TrimSpace(m) == requestModel {
					modelAllowed = true
					break
				}
			}
			if !modelAllowed {
				c.JSON(http.StatusForbidden, gin.H{
					"error": gin.H{
						"message": "Model not allowed: " + requestModel,
						"type":    "invalid_request_error",
					},
				})
				return
			}
		}
	}

	tokenID, _ := c.Get("token_id")
	userID, _ := c.Get("token_user_id")

	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host
			// Keep original path (e.g., /v1/chat/completions)
			req.Header.Set("Authorization", "Bearer "+upstreamKey)
			req.Header.Del("Cookie")
		},
		ModifyResponse: func(resp *http.Response) error {
			duration := int(time.Since(startTime).Milliseconds())

			if isStream {
				// For streaming responses, wrap the body to capture usage
				resp.Body = &streamReader{
					reader:   resp.Body,
					tokenID:  tokenID.(uint),
					userID:   userID.(uint),
					model:    requestModel,
					path:     c.Request.URL.Path,
					method:   c.Request.Method,
					ip:       getRequestIP(c),
					status:   resp.StatusCode,
					duration: duration,
				}
			} else {
				// For non-streaming, read body, extract usage, re-wrap
				body, err := io.ReadAll(resp.Body)
				resp.Body.Close()
				if err == nil {
					var usage UsageInfo
					extractUsageFromJSON(body, &usage)

					logEntry := model.RequestLog{
						UserID:           userID.(uint),
						TokenID:          tokenID.(uint),
						RequestIP:        getRequestIP(c),
						Method:           c.Request.Method,
						Path:             c.Request.URL.Path,
						Model:            requestModel,
						StatusCode:       resp.StatusCode,
						Duration:         duration,
						PromptTokens:     usage.PromptTokens,
						CompletionTokens: usage.CompletionTokens,
						TotalTokens:      usage.TotalTokens,
						CreatedAt:        time.Now(),
					}
					service.RecordLog(logEntry)

					if resp.StatusCode >= 200 && resp.StatusCode < 300 {
						service.IncrementUsage(tokenID.(uint), userID.(uint))
					}

					resp.Body = io.NopCloser(bytes.NewBuffer(body))
					resp.ContentLength = int64(len(body))
				}
			}
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error: %v", err)
			duration := int(time.Since(startTime).Milliseconds())

			logEntry := model.RequestLog{
				UserID:       userID.(uint),
				TokenID:      tokenID.(uint),
				RequestIP:    getRequestIP(c),
				Method:       c.Request.Method,
				Path:         c.Request.URL.Path,
				Model:        requestModel,
				StatusCode:   502,
				Duration:     duration,
				ErrorMessage: err.Error(),
				CreatedAt:    time.Now(),
			}
			service.RecordLog(logEntry)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(gin.H{
				"error": gin.H{
					"message": "Upstream service unavailable",
					"type":    "server_error",
				},
			})
		},
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func getRequestIP(c *gin.Context) string {
	if ip, exists := c.Get("request_ip"); exists {
		return ip.(string)
	}
	return utils.GetClientIP(c)
}

type UsageInfo struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func extractUsageFromJSON(body []byte, usage *UsageInfo) {
	var resp map[string]interface{}
	if json.Unmarshal(body, &resp) != nil {
		return
	}
	if u, ok := resp["usage"].(map[string]interface{}); ok {
		if v, ok := u["prompt_tokens"].(float64); ok {
			usage.PromptTokens = int(v)
		}
		if v, ok := u["completion_tokens"].(float64); ok {
			usage.CompletionTokens = int(v)
		}
		if v, ok := u["total_tokens"].(float64); ok {
			usage.TotalTokens = int(v)
		}
	}
}
