#!/bin/bash

echo "Creating cmd directory if it doesn't exist..."
mkdir -p cmd/test-logs

echo "Building test-logs command line tool..."
go build -o test-logs cmd/test-logs/main.go

if [ $? -eq 0 ]; then
    echo "Build successful! You can now use ./test-logs --auth-key <auth-key>"
    echo "Default auth key is required. Example: AQEtdOfDoeD1vTYmSm4ERnwpPdmVXf0GEKZmGurd1n3RybTVVeIHLB0qo6UvvANUQ-50KvaWxH79zA3-Wweb8ijLOu2BnGnUckIJFx5Y0F_KvJn6B1MojRgtLSPaF_NJW5oBxzqo7g1VVkZ8Nc-1g5z1ro6mbNH8zTqA40KjSWHdyz3ZggXtt_rCCfrpW_Ed6C9qJNXP44TsXX0VV5C0nzhPjJq-uYML4c1Cb3XesV9czXnk4E8rEMNXGjTopBovoSkvTSab-mWi72DSUv9ElgA3EJWDpJN4hsM6oZJjeSR-UIwfSXnLQ8I7gjcL2jpD97hF5nI"
    echo ""
    echo "Additional options:"
    echo "  --duration <seconds>    Duration to run (default: 60)"
    echo "  --destination <url>     Log destination URL (default: https://ingestion.easylogs.co/logs)"
    echo "  --batch-size <count>    Number of logs to send in each batch (default: 10)"
    echo "  --interval <ms>         Interval between batches in milliseconds (default: 1000)"
    echo ""
    echo "This tool will send ALL data types (api, db, user, metrics) without exceptions."
    echo ""
    echo "Example:"
    echo "./test-logs --auth-key YOUR_AUTH_KEY --duration 86400 --batch-size 20 --interval 500"
else
    echo "Build failed!"
fi 