package proxy

import (
	"bufio"
	"cpa-distribution/model"
	"cpa-distribution/service"
	"encoding/json"
	"io"
	"log"
	"strings"
	"time"
)

type streamReader struct {
	reader   io.ReadCloser
	tokenID  uint
	userID   uint
	model    string
	path     string
	method   string
	ip       string
	status   int
	duration int
	usage    UsageInfo
	done     bool
	buffer   []byte
	scanner  *bufio.Scanner
	inited   bool
}

func (s *streamReader) Read(p []byte) (int, error) {
	if !s.inited {
		s.scanner = bufio.NewScanner(s.reader)
		s.scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB buffer
		s.inited = true
	}

	if len(s.buffer) > 0 {
		n := copy(p, s.buffer)
		s.buffer = s.buffer[n:]
		return n, nil
	}

	if s.done {
		return 0, io.EOF
	}

	if s.scanner.Scan() {
		line := s.scanner.Text()

		// Try to extract usage from SSE data
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				s.done = true
				// Record log when stream ends
				s.recordStreamLog()
			} else {
				s.parseStreamChunk(data)
			}
		}

		// Return the line with newline
		output := line + "\n"
		n := copy(p, output)
		if n < len(output) {
			s.buffer = []byte(output[n:])
		}
		return n, nil
	}

	if err := s.scanner.Err(); err != nil {
		return 0, err
	}

	// EOF - record log if not done yet
	if !s.done {
		s.done = true
		s.recordStreamLog()
	}
	return 0, io.EOF
}

func (s *streamReader) parseStreamChunk(data string) {
	var chunk map[string]interface{}
	if json.Unmarshal([]byte(data), &chunk) != nil {
		return
	}

	// Extract usage if present (usually in the last chunk)
	if u, ok := chunk["usage"].(map[string]interface{}); ok {
		if v, ok := u["prompt_tokens"].(float64); ok {
			s.usage.PromptTokens = int(v)
		}
		if v, ok := u["completion_tokens"].(float64); ok {
			s.usage.CompletionTokens = int(v)
		}
		if v, ok := u["total_tokens"].(float64); ok {
			s.usage.TotalTokens = int(v)
		}
	}

	// Also capture model from chunk if not set
	if s.model == "" {
		if m, ok := chunk["model"].(string); ok {
			s.model = m
		}
	}
}

func (s *streamReader) recordStreamLog() {
	logEntry := model.RequestLog{
		UserID:           s.userID,
		TokenID:          s.tokenID,
		RequestIP:        s.ip,
		Method:           s.method,
		Path:             s.path,
		Model:            s.model,
		StatusCode:       s.status,
		Duration:         s.duration,
		PromptTokens:     s.usage.PromptTokens,
		CompletionTokens: s.usage.CompletionTokens,
		TotalTokens:      s.usage.TotalTokens,
		CreatedAt:        time.Now(),
	}
	service.RecordLog(logEntry)

	if s.status >= 200 && s.status < 300 {
		service.IncrementUsage(s.tokenID, s.userID)
	}

	log.Printf("Stream completed: model=%s, tokens=%d, duration=%dms", s.model, s.usage.TotalTokens, s.duration)
}

func (s *streamReader) Close() error {
	if !s.done {
		s.done = true
		s.recordStreamLog()
	}
	return s.reader.Close()
}
