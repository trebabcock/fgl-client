#!/bin/bash

if [ $# -eq 0 ]; then
    echo "Version not provided"
    exit 1
fi

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

rm version.json
echo "{\"version\":\"$1\"}" > version.json
