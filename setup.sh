#!/bin/bash
set -e

echo "DAAR Project 3 - Setup"
echo "======================"


# Download Books
echo "[1/3] Building index..."
go run cmd/download_books/main.go

# Build index
echo "[2/3] Building index..."
go run cmd/indexer/main.go

# Build graph
echo "[3/3] Building graph..."
go run cmd/graph/main.go

echo ""
echo "Done! Run: go run cmd/server/main.go"