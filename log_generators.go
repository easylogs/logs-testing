package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	logLevels     = []string{"INFO", "WARN", "ERROR", "DEBUG"}
	environments  = []string{"production", "staging", "development"}
	apiPaths      = []string{"/api/users", "/api/products", "/api/orders", "/api/auth", "/api/payments"}
	httpMethods   = []string{"GET", "POST", "PUT", "DELETE"}
	userActions   = []string{"login", "logout", "purchase", "view_item", "update_profile"}
	dbOperations  = []string{"SELECT", "INSERT", "UPDATE", "DELETE"}
	services      = []string{"auth-service", "user-service", "payment-service", "inventory-service", "notification-service"}
)

func generateAPILogs(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			batchSize := rand.Intn(4) + 2 // Random number between 2 and 5
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
					Path:       apiPaths[rand.Intn(len(apiPaths))],
					Duration:   duration,
					Environment: environments[rand.Intn(len(environments))],
				}
				broadcastLog(logs[i])
			}
			bulkIndexLogs(logs)
		case <-stopChan:
			return
		}
	}
}

func generateDatabaseLogs(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			batchSize := rand.Intn(4) + 2 // Random number between 2 and 5
			logs := make([]LogEntry, batchSize)
			for i := 0; i < batchSize; i++ {
				operation := dbOperations[rand.Intn(len(dbOperations))]
				duration := rand.Intn(500)

				logs[i] = LogEntry{
					Timestamp:   time.Now().Format(time.RFC3339),
					Level:       logLevels[rand.Intn(len(logLevels))],
					Service:     services[rand.Intn(len(services))],
					Message:     fmt.Sprintf("Database operation %s completed in %dms", operation, duration),
					Duration:    duration,
					Action:      operation,
					Environment: environments[rand.Intn(len(environments))],
					Metadata: map[string]interface{}{
						"query_type": operation,
						"table":      []string{"users", "orders", "products"}[rand.Intn(3)],
					},
				}
				broadcastLog(logs[i])
			}
			bulkIndexLogs(logs)
		case <-stopChan:
			return
		}
	}
}

func generateUserActivityLogs(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			batchSize := rand.Intn(4) + 2 // Random number between 2 and 5
			logs := make([]LogEntry, batchSize)
			for i := 0; i < batchSize; i++ {
				userID := fmt.Sprintf("user_%d", rand.Intn(1000))
				action := userActions[rand.Intn(len(userActions))]

				logs[i] = LogEntry{
					Timestamp:   time.Now().Format(time.RFC3339),
					Level:       "INFO",
					Service:     services[rand.Intn(len(services))],
					Message:     fmt.Sprintf("User %s performed action: %s", userID, action),
					UserID:      userID,
					Action:      action,
					Environment: environments[rand.Intn(len(environments))],
				}
				broadcastLog(logs[i])
			}
			bulkIndexLogs(logs)
		case <-stopChan:
			return
		}
	}
}

func generateSystemMetrics(wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			batchSize := rand.Intn(4) + 2 // Random number between 2 and 5
			logs := make([]LogEntry, batchSize)
			for i := 0; i < batchSize; i++ {
				cpuUsage := rand.Float64() * 100
				memoryUsage := rand.Float64() * 100
				diskUsage := rand.Float64() * 100

				logs[i] = LogEntry{
					Timestamp:   time.Now().Format(time.RFC3339),
					Level:       logLevels[rand.Intn(len(logLevels))],
					Service:     services[rand.Intn(len(services))],
					Message:     fmt.Sprintf("System metrics - CPU: %.2f%%, Memory: %.2f%%, Disk: %.2f%%", 
						cpuUsage, memoryUsage, diskUsage),
					Environment: environments[rand.Intn(len(environments))],
					Metadata: map[string]interface{}{
						"cpu_usage":    cpuUsage,
						"memory_usage": memoryUsage,
						"disk_usage":   diskUsage,
					},
				}
				broadcastLog(logs[i])
			}
			bulkIndexLogs(logs)
		case <-stopChan:
			return
		}
	}
} 