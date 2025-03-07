package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"time"
)

const (
	// Old configurations (commented out)
	// logzioURL = "http://listener.logz.io:8070"
	// logzioToken = "zFFZfgoEffBEjcvCneBfOqSNpkCWQdJE"  // Your Logz.io token
	
	// Old OpenSearch configuration
	// elasticHost = "https://5.161.36.161:443"
	// elasticIndex = "logs-6-23368a714a619f2a"
	// authHeader = "Basic dXNlcl91b0FIZ3FVSzBhOXhjamRROnNrMHRMRVpJT0V5TVRJRzc2Y0Y0MHRGcWRUWXc3Zlho"
	
	// EasyLogs configuration
	elasticHost = "https://ingestion.easylogs.co/logs"
	authHeader = "Bearer AQEtdOfDoeD1vTYmSm4ERnwpPdmVXf0GEKZmGurd1n3RybTVVeIHLB0qo6UvvANUQ-50KvaWxH79zA3-Wweb8ijLOu2BnGnUckIJFx5Y0F_KvJn6B1MojRgtLSPaF_NJW5oBxzqo7g1VVkZ8Nc-1g5z1ro6mbNH8zTqA40KjSWHdyz3ZggXtt_rCCfrpW_Ed6C9qJNXP44TsXX0VV5C0nzhPjJq-uYML4c1Cb3XesV9czXnk4E8rEMNXGjTopBovoSkvTSab-mWi72DSUv9ElgA3EJWDpJN4hsM6oZJjeSR-UIwfSXnLQ8I7gjcL2jpD97hF5nI"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp   string      `json:"timestamp"`
	Level       string      `json:"level"`
	Service     string      `json:"service"`
	Message     string      `json:"message"`
	StatusCode  int         `json:"status_code,omitempty"`
	Method      string      `json:"method,omitempty"`
	Path        string      `json:"path,omitempty"`
	Duration    int         `json:"duration,omitempty"`
	UserID      string      `json:"user_id,omitempty"`
	Action      string      `json:"action,omitempty"`
	Metadata    interface{} `json:"metadata,omitempty"`
	Environment string      `json:"environment"`
}

func bulkIndexLogs(logs []LogEntry) {
	// Create a buffer for the JSON array of logs
	var buf bytes.Buffer
	
	// Encode the entire logs array as a single JSON array
	if err := json.NewEncoder(&buf).Encode(logs); err != nil {
		stdlog.Printf("Error encoding log entries: %s", err)
		return
	}
	
	// Create HTTP request to EasyLogs
	req, err := http.NewRequest("POST", elasticHost, &buf)
	if err != nil {
		stdlog.Printf("Error creating request: %s", err)
		return
	}
	
	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// Send request
	resp, err := client.Do(req)
	if err != nil {
		stdlog.Printf("Error sending logs: %s", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		stdlog.Printf("Error from API. Status: %d", resp.StatusCode)
		// Read and log response body for debugging
		body, _ := ioutil.ReadAll(resp.Body)
		stdlog.Printf("Response: %s", string(body))
	}
} 