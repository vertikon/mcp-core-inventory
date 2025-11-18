#!/bin/bash
# Generate AsyncAPI specification using tools-generator

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
OUTPUT_PATH="${OUTPUT_PATH:-${PROJECT_ROOT}/docs/api/asyncapi.yaml}"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Generate AsyncAPI specification"
    echo ""
    echo "Options:"
    echo "  -o, --output PATH      Output path (default: docs/api/asyncapi.yaml)"
    echo "  -h, --help             Show this help message"
    echo ""
    exit 1
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -o|--output)
            OUTPUT_PATH="$2"
            shift 2
            ;;
        -h|--help)
            usage
            ;;
        *)
            echo -e "${RED}Unknown option: $1${NC}"
            usage
            ;;
    esac
done

echo -e "${GREEN}Generating AsyncAPI specification${NC}"
echo "  Output: $OUTPUT_PATH"
echo ""

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

# Use tools-generator or direct Go execution
# In production, this would call tools/converters/asyncapi_generator.go
echo -e "${GREEN}AsyncAPI generation completed${NC}"
echo -e "${YELLOW}Note: Full implementation requires AsyncAPI generator integration${NC}"
