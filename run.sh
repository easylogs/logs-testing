#!/bin/bash

echo "Starting the log generator application..."
go run *.go

# This will run the main application which starts the web server
# You can then use the test-logs command line tool to interact with it 