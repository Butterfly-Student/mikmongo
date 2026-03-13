#!/bin/bash

# Script untuk menambahkan router Mikrotik via API
# Usage: ./add_router.sh [BASE_URL]

BASE_URL=${1:-"http://localhost:8080"}

echo "Menambahkan router ke API: $BASE_URL"
echo ""

curl -X POST "$BASE_URL/api/v1/routers" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Testing",
    "address": "192.168.233.1",
    "api_port": 8728,
    "rest_port": 80,
    "username": "admin",
    "password": "r00t",
    "is_master": false,
    "is_active": true,
    "status": "online"
  }' \
  -w "\n\nHTTP Status: %{http_code}\n" \
  -v

echo ""
echo "Selesai!"
