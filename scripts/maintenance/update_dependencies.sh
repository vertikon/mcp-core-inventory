#!/bin/bash
# Update dependencies

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Update dependencies"
    echo ""
    echo "Options:"
    echo "  -h, --help             Show this help message"
    echo ""
    exit 1
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            usage
            ;;
    esac
done

echo -e "${GREEN}Updating dependencies...${NC}"
echo ""

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

# Update Go dependencies
echo -e "${GREEN}Updating Go dependencies...${NC}"
cd "$PROJECT_ROOT"
go get -u ./... || {
    echo -e "${YELLOW}Warning: Some dependencies may have failed to update${NC}"
}

go mod tidy

echo -e "${GREEN}Dependencies update completed${NC}"
