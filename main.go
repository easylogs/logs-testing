package main

import (
	"sync"
)

// Configuration constants
const (
	batchSize     = 100
)

var (
	// Control channel for stopping log generation
	stopChan = make(chan struct{})
	// Track if generators are running
	isRunning bool
	// Mutex for controlling isRunning access
	runningMux sync.Mutex
)

func main() {
	// Start the web server
	startWebServer()
}

// Function to stop log generation
func stopLogGeneration() {
	runningMux.Lock()
	defer runningMux.Unlock()
	
	if isRunning {
		close(stopChan)
		isRunning = false
		// Create new channel for next start
		stopChan = make(chan struct{})
	}
} 