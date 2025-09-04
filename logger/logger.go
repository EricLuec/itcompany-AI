package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var logAPI = "http://localhost:8080/logEntry"

type LogEntry struct {
	Category    string `json:"category"`
	Description string `json:"description"`
	Timestamp   string `json:"timestamp"`
}

func CreateLogEntry(category, description string) error {
	logEntry := LogEntry{
		Category:    category,
		Description: description,
		Timestamp:   time.Now().Format("2006-01-02T15:04:05"),
	}

	logData, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("error marshaling log entry: %v", err)
	}

	resp, err := http.Post(logAPI, "application/json", bytes.NewBuffer(logData))
	if err != nil {
		return fmt.Errorf("error posting log entry: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with: %s", resp.Status)
	}

	log.Println("Log entry posted successfully.")
	return nil
}
