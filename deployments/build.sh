#!/bin/bash

APP_NAME="api-server"

rm -f $APP_NAME

echo "Compile $APP_NAME for Linux amd64."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ../cmd/api-server/

if [ $? -eq 0 ]; then
    echo "✅ success."
else
    echo "❌ failed."
    exit 1
fi
