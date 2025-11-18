#!/bin/bash
# Validate MCP project using tools-validator

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
MCP_PATH="${MCP_PATH:-}"
STRICT_MODE="${STRICT_MODE:-false}"
CHECK_SECURITY="${CHECK_SECURITY:-false}"
CHECK_DEPENDENCIES="${CHECK_DEPENDENCIES:-false}"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Validate MCP project using tools-validator"
    echo ""
    echo "Options:"
    echo "  -p, --path PATH        Path to MCP project (required)"
    echo "  -s, --strict           Enable strict mode"
    echo "  --security             Check security"
    echo "  --dependencies          Check dependencies"
    echo "  -h, --help             Show this help message"
    echo ""
    exit 1
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--path)
            MCP_PATH="$2"
            shift 2
            ;;
        -s|--strict)
            STRICT_MODE="true"
            shift
            ;;
        --security)
            CHECK_SECURITY="true"
            shift
            ;;
        --dependencies)
            CHECK_DEPENDENCIES="true"
            shift
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

# Validate required parameters
if [ -z "$MCP_PATH" ]; then
    echo -e "${RED}Error: MCP path is required${NC}"
    usage
fi

# Resolve absolute path
MCP_PATH="$(cd "$MCP_PATH" && pwd)"

echo -e "${GREEN}Validating MCP project${NC}"
echo "  Path: $MCP_PATH"
echo "  Strict mode: $STRICT_MODE"
echo "  Check security: $CHECK_SECURITY"
echo "  Check dependencies: $CHECK_DEPENDENCIES"
echo ""

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

# Build tools-validator if needed
TOOLS_VALIDATOR="${PROJECT_ROOT}/bin/tools-validator"
if [ ! -f "$TOOLS_VALIDATOR" ]; then
    echo -e "${YELLOW}Building tools-validator...${NC}"
    cd "$PROJECT_ROOT"
    mkdir -p bin
    go build -o "$TOOLS_VALIDATOR" ./cmd/tools-validator || {
        echo -e "${RED}Error: Failed to build tools-validator${NC}"
        exit 1
    }
fi

# Execute validation
echo -e "${GREEN}Executing MCP validation...${NC}"
cd "$PROJECT_ROOT"

# Build command
CMD="$TOOLS_VALIDATOR -type mcp -path \"$MCP_PATH\""
[ "$STRICT_MODE" = "true" ] && CMD="$CMD -strict"
[ "$CHECK_SECURITY" = "true" ] && CMD="$CMD -security"
[ "$CHECK_DEPENDENCIES" = "true" ] && CMD="$CMD -dependencies"

# Execute
eval $CMD
VALIDATION_RESULT=$?

if [ $VALIDATION_RESULT -eq 0 ]; then
    echo -e "${GREEN}MCP validation passed${NC}"
    exit 0
else
    echo -e "${RED}Error: MCP validation failed${NC}"
    exit 1
fi
