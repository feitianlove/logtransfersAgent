#!/bin/bash
dos2unix $0 > /dev/null 2>&1

cd "$(dirname "$0")"

#test -f tflow && rm tflow

echo "build..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.branches=$(git symbolic-ref --short -q HEAD) -X main.gitRev=$(git rev-parse HEAD) -X main.buildTime=$(date +'%Y-%m-%d_%T')" -o logtransferAgent ../main.go
echo "build success"