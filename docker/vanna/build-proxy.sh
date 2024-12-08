#!/bin/bash

go fmt ./...
go vet ./...
go mod tidy

go build -o local/bin/proxy main.go
