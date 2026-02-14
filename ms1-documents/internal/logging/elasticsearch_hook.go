package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap/zapcore"
)

type ElasticsearchHook struct {
	client      *http.Client
	esURL       string
	indexPrefix string
}

func NewElasticsearchHook() *ElasticsearchHook {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://localhost:9200"
	}

	return &ElasticsearchHook{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		esURL:       esURL,
		indexPrefix: "ms1-logs",
	}
}

func (h *ElasticsearchHook) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	logDoc := map[string]interface{}{
		"@timestamp": entry.Time.Format(time.RFC3339),
		"level":      entry.Level.String(),
		"message":    entry.Message,
		"service":    "ms1-documents",
		"caller":     entry.Caller.String(),
	}

	for _, field := range fields {
		logDoc[field.Key] = field.Interface
	}

	go h.sendToElasticsearch(logDoc)

	return nil
}

func (h *ElasticsearchHook) sendToElasticsearch(logDoc map[string]interface{}) {
	indexName := fmt.Sprintf("%s-%s", h.indexPrefix, time.Now().Format("2006.01.02"))
	url := fmt.Sprintf("%s/%s/_doc", h.esURL, indexName)

	data, err := json.Marshal(logDoc)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
