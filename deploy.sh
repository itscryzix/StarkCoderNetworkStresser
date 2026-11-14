#!/bin/bash
# Stark Coder Public Deployment Script

echo "=== Stark Coder Deployment ==="
echo ""

# Check prerequisites
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed"
    exit 1
fi

if ! command -v node &> /dev/null; then
    echo "ERROR: Node.js is not installed"
    exit 1
fi

echo "[1/4] Installing frontend dependencies..."
cd web-client
npm install --production
cd ..

echo "[2/4] Building backend..."
go build -o bin/mmb-server ./cmd/server

echo "[3/4] Building frontend..."
cd web-client
npm run build
cd ..
cp -r web-client/dist/* bin/web-client/

echo "[4/4] Cleaning up..."
rm -rf node_modules/.cache

echo ""
echo "=== Build Complete ==="
echo "Server: http://localhost:3000"
echo "Binary: bin/mmb-server"
echo ""
echo "To deploy publicly, see DEPLOY.md"
