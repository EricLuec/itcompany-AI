package logger

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateLogEntry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var entry LogEntry
		if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
			t.Errorf("decode error: %v", err)
		}
		if entry.Category != "test" || entry.Description != "testing" {
			t.Errorf("unexpected log entry: %+v", entry)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	origURL := logAPI
	logAPI = server.URL
	defer func() { logAPI = origURL }()

	err := CreateLogEntry("test", "testing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
