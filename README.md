# Logs Testing Application

This application is designed for testing log collection and forwarding to various log management services. It provides a flexible way to generate and send logs to different destinations for testing and development purposes.

## Overview

The application can be configured to send logs to different log management services by modifying the configuration in `opensearch_helpers.go`. Currently, it's set up to work with EasyLogs, but it can be reconfigured to work with other services like Logz.io or OpenSearch.

## Configuration

To change the log destination, edit the constants in `opensearch_helpers.go` (around line 12-25):

```go
const (
    // Configure your log destination here
    elasticHost = "https://ingestion.easylogs.co/logs"
    authHeader = "Bearer YOUR_AUTH_TOKEN"
)
```

## Available Log Destinations

The code includes commented examples for different log destinations:

1. **Logz.io**:
   ```go
   logzioURL = "http://listener.logz.io:8070"
   logzioToken = "YOUR_LOGZ_IO_TOKEN"
   ```

2. **OpenSearch**:
   ```go
   elasticHost = "https://your-opensearch-host:443"
   elasticIndex = "your-index-name"
   authHeader = "Basic YOUR_BASE64_ENCODED_CREDENTIALS"
   ```

3. **EasyLogs** (current configuration):
   ```go
   elasticHost = "https://ingestion.easylogs.co/logs"
   authHeader = "Bearer YOUR_AUTH_TOKEN"
   ```

## Log Structure

The application uses a structured log format defined in the `LogEntry` struct, which includes fields like:
- Timestamp
- Log level
- Service name
- Message
- HTTP status code
- Method
- Path
- Duration
- User ID
- Action
- Metadata
- Environment

## Usage

### Web Interface

Run the application to generate and send test logs to your configured destination:

```
go run *.go
```

Then open your browser to http://localhost:8090 to access the web interface.

### Command Line Tool

A standalone command line tool is available for testing log generation. This tool sends logs directly to your log destination without requiring the web server to be running:

1. Build the tool:
   ```
   ./build.sh
   ```

2. Run the tool with an authentication key:
   ```
   ./test-logs --auth-key YOUR_AUTH_KEY
   ```

#### Command Line Options

- `--auth-key <key>`: Authentication key for the log destination (required)
- `--duration <seconds>`: Duration to run log generation (default: 60 seconds)
- `--destination <url>`: Log destination URL (default: https://ingestion.easylogs.co/logs)
- `--batch-size <count>`: Number of logs to send in each batch (default: 10)
- `--interval <ms>`: Interval between batches in milliseconds (default: 1000)

The command line tool will send ALL data types (api, db, user, metrics) without exceptions.

Example:
```
./test-logs --auth-key YOUR_AUTH_KEY --duration 86400 --batch-size 20 --interval 500
```

This will run the log generator for 24 hours, sending batches of 20 logs every 500 milliseconds directly to the log destination.

## Note

Remember to replace authentication tokens and credentials with your own when configuring different log destinations. 