#!/bin/bash

set -e

# Configuration
GATEWAY_URL="${GATEWAY_URL:-http://localhost:8080}"
AUTH0_TOKEN="${AUTH0_TOKEN}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper function to print colored output
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Check if token is provided
if [ -z "$AUTH0_TOKEN" ]; then
    print_error "AUTH0_TOKEN environment variable is not set"
    echo "Usage: AUTH0_TOKEN=your_token_here $0"
    exit 1
fi

print_info "Testing Gateway API (gateway.v1.MeService) at $GATEWAY_URL"
echo ""

# Test: GetMe - Get current authenticated user information
print_info "Test: GetMe (gateway.v1.MeService/GetMe)"
RESPONSE=$(curl -s -w "\n%{http_code}" \
    -H "Authorization: Bearer $AUTH0_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{}' \
    "$GATEWAY_URL/gateway.v1.MeService/GetMe")

HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)
BODY=$(echo "$RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    print_success "GetMe succeeded"
    echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
else
    print_error "GetMe failed with status $HTTP_CODE"
    echo "$BODY"
fi
echo ""

print_info "Testing complete"
