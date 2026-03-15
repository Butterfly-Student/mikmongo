#!/bin/bash

# Test all mikhmon scenarios automatically
# Usage: ./test_all.sh

cd "$(dirname "$0")"

echo "====================================="
echo "  Mikhmon Complete Test Suite"
echo "  Router: 192.168.233.1:8728"
echo "====================================="
echo ""

# Function to run a test scenario
run_test() {
    local scenario=$1
    echo ""
    echo "====================================="
    echo "  Testing Scenario $scenario"
    echo "====================================="
    echo "$scenario" | timeout 60 go run main.go 2>&1 | head -100
}

# Test 1: Voucher Generator
echo "[1/5] Testing Voucher Generator..."
run_test "1"

# Test 2: Profile Manager  
echo ""
echo "[2/5] Testing Profile Manager..."
run_test "2"

# Test 3: Multi Router
echo ""
echo "[3/5] Testing Multi Router..."
run_test "3"

# Test 4: Report Viewer
echo ""
echo "[4/5] Testing Report Viewer..."
run_test "4"

# Test 5: Expire Monitor
echo ""
echo "[5/5] Testing Expire Monitor..."
run_test "5"

echo ""
echo "====================================="
echo "  All Tests Completed!"
echo "====================================="
