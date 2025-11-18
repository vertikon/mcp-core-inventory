#!/bin/bash
# Generate MCP project using tools-generator

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Load configuration
CONFIG_DIR="${PROJECT_ROOT}/config"
ENV="${HULK_ENV:-dev}"

# Default values
MCP_NAME="${MCP_NAME:-}"
OUTPUT_PATH="${OUTPUT_PATH:-}"
STACK="${STACK:-mcp-go-premium}"
FEATURES="${FEATURES:-}"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Generate MCP project using tools-generator"
    echo ""
    echo "Options:"
    echo "  -n, --name NAME        MCP project name (required)"
    echo "  -o, --output PATH      Output path (required)"
    echo "  -s, --stack STACK      Stack type (default: mcp-go-premium)"
    echo "  -f, --features LIST    Comma-separated list of features"
    echo "  -h, --help             Show this help message"
    echo ""
    echo "Environment variables:"
    echo "  MCP_NAME              MCP project name"
    echo "  OUTPUT_PATH            Output path"
    echo "  STACK                  Stack type"
    echo "  FEATURES               Comma-separated features"
    exit 1
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--name)
            MCP_NAME="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_PATH="$2"
            shift 2
            ;;
        -s|--stack)
            STACK="$2"
            shift 2
            ;;
        -f|--features)
            FEATURES="$2"
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

# Validate required parameters
if [ -z "$MCP_NAME" ]; then
    echo -e "${RED}Error: MCP name is required${NC}"
    usage
fi

if [ -z "$OUTPUT_PATH" ]; then
    echo -e "${RED}Error: Output path is required${NC}"
    usage
fi

echo -e "${GREEN}Generating MCP project${NC}"
echo "  Name: $MCP_NAME"
echo "  Output: $OUTPUT_PATH"
echo "  Stack: $STACK"
[ -n "$FEATURES" ] && echo "  Features: $FEATURES"
echo ""

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

# Build tools-generator if needed
TOOLS_GENERATOR="${PROJECT_ROOT}/bin/tools-generator"
if [ ! -f "$TOOLS_GENERATOR" ]; then
    echo -e "${YELLOW}Building tools-generator...${NC}"
    cd "$PROJECT_ROOT"
    mkdir -p bin
    go build -o "$TOOLS_GENERATOR" ./cmd/tools-generator || {
        echo -e "${RED}Error: Failed to build tools-generator${NC}"
        exit 1
    }
fi

# Execute generation
echo -e "${GREEN}Executing MCP generation...${NC}"
cd "$PROJECT_ROOT"

# Build command
CMD="$TOOLS_GENERATOR -type mcp -name \"$MCP_NAME\" -path \"$OUTPUT_PATH\" -stack \"$STACK\""
[ -n "$FEATURES" ] && CMD="$CMD -features \"$FEATURES\""

# Execute
eval $CMD

if [ $? -eq 0 ]; then
    echo -e "${GREEN}MCP generation completed successfully${NC}"
else
    echo -e "${RED}Error: MCP generation failed${NC}"
    exit 1
fi
