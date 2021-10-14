#! /bin/bash

go build cmd/csvlog/main.go -o csvlog_linux_amd64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o csvlog_win_amd64 cmd/csvlog/main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o csvlog_darwin_amd64 cmd/csvlog/main.go
