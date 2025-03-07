package main

import (
	"fmt"
	"html/template"
	stdlog "log"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Channel to broadcast logs to all connected clients
	logBroadcast = make(chan LogEntry, 100)
	// Active WebSocket connections
	clients    = make(map[*websocket.Conn]bool)
	clientsMux sync.Mutex
)

// Theme defines the color scheme and styling
type Theme struct {
	BackgroundColor   string
	PrimaryColor      string
	PrimaryColorHover string
	InfoColor         string
	WarnColor         string
	ErrorColor        string
	DebugColor        string
	DangerColor      string
	DangerColorHover string
}

// PageData contains all the data needed for the template
type PageData struct {
	Title               string
	Header              string
	ButtonText          string
	ButtonTextAfterStart string
	MaxLogEntries       int
	Theme               Theme
}

func getDefaultPageData() PageData {
	return PageData{
		Title:               "Log Generator Dashboard",
		Header:              "Real-time Log Generator Dashboard",
		ButtonText:          "Start Log Generation",
		ButtonTextAfterStart: "Log Generation Started",
		MaxLogEntries:       1000,
		Theme: Theme{
			BackgroundColor:   "#f0f0f0",
			PrimaryColor:      "#4CAF50",
			PrimaryColorHover: "#45a049",
			InfoColor:         "#2196F3",
			WarnColor:         "#FF9800",
			ErrorColor:        "#f44336",
			DebugColor:        "#4CAF50",
			DangerColor:      "#f44336",
			DangerColorHover: "#d32f2f",
		},
	}
}

func startWebServer() {
	// Serve static files
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/stop", handleStop)

	fmt.Println("Web server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		stdlog.Fatal(err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", "index.html"))
	if err != nil {
		stdlog.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := getDefaultPageData()
	err = tmpl.Execute(w, data)
	if err != nil {
		stdlog.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		stdlog.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Register new client
	clientsMux.Lock()
	clients[conn] = true
	clientsMux.Unlock()

	// Clean up on disconnect
	defer func() {
		clientsMux.Lock()
		delete(clients, conn)
		clientsMux.Unlock()
		conn.Close()
	}()

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	runningMux.Lock()
	if isRunning {
		runningMux.Unlock()
		http.Error(w, "Log generation already running", http.StatusConflict)
		return
	}
	// Create new stop channel
	stopChan = make(chan struct{})
	isRunning = true
	runningMux.Unlock()

	// Start the log generators
	var wg sync.WaitGroup
	wg.Add(4)

	go generateAPILogs(&wg)
	go generateDatabaseLogs(&wg)
	go generateUserActivityLogs(&wg)
	go generateSystemMetrics(&wg)

	// Start a goroutine to wait for completion
	go func() {
		wg.Wait()
		runningMux.Lock()
		isRunning = false
		runningMux.Unlock()
	}()

	w.Write([]byte("Log generation started"))
}

func handleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	runningMux.Lock()
	if !isRunning {
		runningMux.Unlock()
		http.Error(w, "Log generation not running", http.StatusBadRequest)
		return
	}
	// Close the stop channel
	close(stopChan)
	isRunning = false
	runningMux.Unlock()

	// Wait a bit longer to ensure all goroutines have stopped
	time.Sleep(500 * time.Millisecond)
	w.Write([]byte("Log generation stopped"))
}

func broadcastLog(log LogEntry) {
	clientsMux.Lock()
	defer clientsMux.Unlock()

	for client := range clients {
		err := client.WriteJSON(log)
		if err != nil {
			stdlog.Printf("Error broadcasting to client: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}