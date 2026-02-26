package service

import (
	"cpa-distribution/model"
	"log"
	"time"
)

var logChannel chan model.RequestLog

func InitLogService() {
	logChannel = make(chan model.RequestLog, 1000)
	go logConsumer()
}

func RecordLog(logEntry model.RequestLog) {
	select {
	case logChannel <- logEntry:
	default:
		log.Println("Log channel full, dropping log entry")
	}
}

func logConsumer() {
	buffer := make([]model.RequestLog, 0, 50)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case entry := <-logChannel:
			buffer = append(buffer, entry)
			if len(buffer) >= 50 {
				flushLogs(buffer)
				buffer = make([]model.RequestLog, 0, 50)
			}
		case <-ticker.C:
			if len(buffer) > 0 {
				flushLogs(buffer)
				buffer = make([]model.RequestLog, 0, 50)
			}
		}
	}
}

func flushLogs(logs []model.RequestLog) {
	if err := model.BatchInsertLogs(logs); err != nil {
		log.Printf("Failed to flush %d logs: %v", len(logs), err)
	}
}

func CleanupOldLogs(days int) (int64, error) {
	return model.DeleteLogsBeforeDays(days)
}
