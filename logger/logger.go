package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type LogEntry struct {
	ID          int       `json:"id"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}

func CreateLogEntry(category, description string, id int) error {
	logEntry := LogEntry{
		ID:          id,
		Category:    category,
		Description: description,
		Timestamp:   time.Now(),
	}

	logEntryJSON, err := json.Marshal(logEntry)
	if err != nil {
		return fmt.Errorf("error marshalling log entry: %v", err)
	}

	url := "http://localhost:8080/logEntry"
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(logEntryJSON))
	if err != nil {
		return fmt.Errorf("error posting log entry: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to log entry, server responded with status: %s", resp.Status)
	}

	log.Printf("Log entry created for category %s: %s", category, description)
	return nil
}
