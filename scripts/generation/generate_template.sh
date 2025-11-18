#!/bin/bash
# Generate template project using tools-generator

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
TEMPLATE_NAME="${TEMPLATE_NAME:-}"
PROJECT_NAME="${PROJECT_NAME:-}"
OUTPUT_PATH="${OUTPUT_PATH:-}"
FEATURES="${FEATURES:-}"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Generate template project using tools-generator"
    echo ""
    echo "Options:"
    echo "  -t, --template NAME   Template name (base/go/tinygo/wasm/web) (required)"
    echo "  -n, --name NAME        Project name (required)"
    echo "  -o, --output PATH     Output path (required)"
    echo "  -f, --features LIST   Comma-separated list of features"
    echo "  -h, --help             Show this help message"
    echo ""
    exit 1
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--template)
            TEMPLATE_NAME="$2"
            shift 2
            ;;
        -n|--name)
            PROJECT_NAME="$2"
            shift 2
            ;;
        -o|--output)
            OUTPUT_PATH="$2"
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
if [ -z "$TEMPLATE_NAME" ]; then
    echo -e "${RED}Error: Template name is required${NC}"
    usage
fi

if [ -z "$PROJECT_NAME" ]; then
    echo -e "${RED}Error: Project name is required${NC}"
    usage
fi

if [ -z "$OUTPUT_PATH" ]; then
    echo -e "${RED}Error: Output path is required${NC}"
    usage
fi

echo -e "${GREEN}Generating template project${NC}"
echo "  Template: $TEMPLATE_NAME"
echo "  Project: $PROJECT_NAME"
echo "  Output: $OUTPUT_PATH"
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
echo -e "${GREEN}Executing template generation...${NC}"
cd "$PROJECT_ROOT"

# Build command
CMD="$TOOLS_GENERATOR -type template -name \"$PROJECT_NAME\" -path \"$OUTPUT_PATH\" -stack \"$TEMPLATE_NAME\""
[ -n "$FEATURES" ] && CMD="$CMD -features \"$FEATURES\""

# Execute
eval $CMD

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Template generation completed successfully${NC}"
else
    echo -e "${RED}Error: Template generation failed${NC}"
    exit 1
fi
