package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Configuration
var (
	authKey    string
	duration   int
	destination string
	batchSize  int
	interval   int
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

// Sample data for log generation
var (
	logLevels     = []string{"INFO", "WARN", "ERROR", "DEBUG"}
	environments  = []string{"production", "staging", "development"}
	apiPaths      = []string{"/api/users", "/api/products", "/api/orders", "/api/auth", "/api/payments"}
	httpMethods   = []string{"GET", "POST", "PUT", "DELETE"}
	userActions   = []string{"login", "logout", "purchase", "view_item", "update_profile"}
	dbOperations  = []string{"SELECT", "INSERT", "UPDATE", "DELETE"}
	services      = []string{"auth-service", "user-service", "payment-service", "inventory-service", "notification-service"}
	userIDs       = []string{"user123", "user456", "user789", "user101", "user202"}
)

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	// Parse command line flags
	flag.StringVar(&authKey, "auth-key", "", "Authentication key for the log destination")
	flag.IntVar(&duration, "duration", 60, "Duration in seconds to run the log generator")
	flag.StringVar(&destination, "destination", "https://ingestion.easylogs.co/logs", "Log destination URL")
	flag.IntVar(&batchSize, "batch-size", 10, "Number of logs to send in each batch")
	flag.IntVar(&interval, "interval", 1000, "Interval between batches in milliseconds")
	flag.Parse()

	// Validate auth key
	if authKey == "" {
		fmt.Println("Error: Authentication key is required")
		flag.Usage()
		os.Exit(1)
	}

	// Start log generation
	fmt.Printf("Starting log generation with auth key: %s\n", authKey)
	fmt.Printf("Duration: %d seconds\n", duration)
	fmt.Printf("Destination: %s\n", destination)
	fmt.Printf("Batch size: %d logs\n", batchSize)
	fmt.Printf("Interval: %d ms\n", interval)
	fmt.Println("Sending ALL data types (api, db, user, metrics)")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Create a stop channel
	stopChan := make(chan struct{})
	
	// Create a wait group for the generators
	var wg sync.WaitGroup
	wg.Add(4)
	
	// Start the log generators
	go generateAPILogs(&wg, stopChan)
	go generateDatabaseLogs(&wg, stopChan)
	go generateUserActivityLogs(&wg, stopChan)
	go generateSystemMetrics(&wg, stopChan)
	
	// Create a timer for the duration
	timer := time.NewTimer(time.Duration(duration) * time.Second)
	
	// Wait for either duration to expire or signal
	select {
	case <-timer.C:
		fmt.Println("Duration completed, stopping log generation...")
		close(stopChan)
	case sig := <-sigChan:
		fmt.Printf("Received signal %v, stopping log generation...\n", sig)
		close(stopChan)
	}
	
	// Wait for all generators to complete
	wg.Wait()
	
	fmt.Println("Log generation stopped successfully")
}

// Send logs to the destination
func sendLogs(logs []LogEntry) {
	// Create a buffer for the JSON array of logs
	var buf bytes.Buffer
	
	// Encode the entire logs array as a single JSON array
	if err := json.NewEncoder(&buf).Encode(logs); err != nil {
		fmt.Printf("Error encoding log entries: %s\n", err)
		return
	}
	
	// Create HTTP request to the destination
	req, err := http.NewRequest("POST", destination, &buf)
	if err != nil {
		fmt.Printf("Error creating request: %s\n", err)
		return
	}
	
	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authKey)
	
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// Send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending logs: %s\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode >= 400 {
		fmt.Printf("Error from API. Status: %d\n", resp.StatusCode)
		// Read and log response body for debugging
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Response: %s\n", string(body))
	} else {
		fmt.Printf("Successfully sent %d logs\n", len(logs))
	}
}

// Generate API logs
func generateAPILogs(wg *sync.WaitGroup, stopChan <-chan struct{}) {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logs := make([]LogEntry, batchSize)
			for i := 0; i < batchSize; i++ {
				statusCode := []int{200, 201, 400, 401, 403, 404, 500}[rand.Intn(7)]
				duration := rand.Intn(1000)

				logs[i] = LogEntry{
					Timestamp:   time.Now().Format(time.RFC3339),
					Level:       logLevels[rand.Intn(len(logLevels))],
					Service:     services[rand.Intn(len(services))],
					Message:     fmt.Sprintf("HTTP %s %s completed in %dms with status %d", 
						httpMethods[rand.Intn(len(httpMethods))],
						apiPaths[rand.Intn(len(apiPaths))],
						duration,
						statusCode),
					StatusCode:  statusCode,
					Method:      httpMethods[rand.Intn(len(httpMethods))],
					Path:        apiPaths[rand.Intn(len(apiPaths))],
					Duration:    duration,
					Environment: environments[rand.Intn(len(environments))],
				}
			}
			sendLogs(logs)
		case <-stopChan:
			return
		}
	}
}

// Generate database logs
func generateDatabaseLogs(wg *sync.WaitGroup, stopChan <-chan struct{}) {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logs := make([]LogEntry, batchSize)
			for i := 0; i < batchSize; i++ {
				duration := rand.Intn(500)
				operation := dbOperations[rand.Intn(len(dbOperations))]
				table := []string{"users", "products", "orders", "payments", "inventory"}[rand.Intn(5)]

				logs[i] = LogEntry{
					Timestamp:   time.Now().Format(time.RFC3339),
					Level:       logLevels[rand.Intn(len(logLevels))],
					Service:     services[rand.Intn(len(services))],
					Message:     fmt.Sprintf("Database operation %s on table %s completed in %dms", 
						operation, table, duration),
					Duration:    duration,
					Method:      operation,
					Path:        table,
					Environment: environments[rand.Intn(len(environments))],
					Metadata:    map[string]interface{}{
						"table":     table,
						"operation": operation,
						"rows":      rand.Intn(100) + 1,
					},
				}
			}
			sendLogs(logs)
		case <-stopChan:
			return
		}
	}
}

// Generate user activity logs
func generateUserActivityLogs(wg *sync.WaitGroup, stopChan <-chan struct{}) {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logs := make([]LogEntry, batchSize)
			for i := 0; i < batchSize; i++ {
				userID := userIDs[rand.Intn(len(userIDs))]
				action := userActions[rand.Intn(len(userActions))]

				logs[i] = LogEntry{
					Timestamp:   time.Now().Format(time.RFC3339),
					Level:       "INFO",
					Service:     "user-activity-service",
					Message:     fmt.Sprintf("User %s performed action: %s", userID, action),
					UserID:      userID,
					Action:      action,
					Environment: environments[rand.Intn(len(environments))],
					Metadata:    map[string]interface{}{
						"browser":  []string{"Chrome", "Firefox", "Safari", "Edge"}[rand.Intn(4)],
						"platform": []string{"Windows", "MacOS", "Linux", "iOS", "Android"}[rand.Intn(5)],
						"ip":       fmt.Sprintf("192.168.%d.%d", rand.Intn(255), rand.Intn(255)),
					},
				}
			}
			sendLogs(logs)
		case <-stopChan:
			return
		}
	}
}

// Generate system metrics
func generateSystemMetrics(wg *sync.WaitGroup, stopChan <-chan struct{}) {
	defer wg.Done()
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			logs := make([]LogEntry, batchSize)
			for i := 0; i < batchSize; i++ {
				cpuUsage := rand.Float64() * 100
				memoryUsage := rand.Float64() * 100
				diskUsage := rand.Float64() * 100

				logs[i] = LogEntry{
					Timestamp:   time.Now().Format(time.RFC3339),
					Level:       "INFO",
					Service:     "system-metrics",
					Message:     fmt.Sprintf("System metrics: CPU: %.2f%%, Memory: %.2f%%, Disk: %.2f%%", 
						cpuUsage, memoryUsage, diskUsage),
					Environment: environments[rand.Intn(len(environments))],
					Metadata:    map[string]interface{}{
						"cpu":    cpuUsage,
						"memory": memoryUsage,
						"disk":   diskUsage,
						"host":   fmt.Sprintf("server-%d", rand.Intn(10)+1),
					},
				}
			}
			sendLogs(logs)
		case <-stopChan:
			return
		}
	}
} 