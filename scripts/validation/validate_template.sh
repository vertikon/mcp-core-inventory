#!/bin/bash
# Validate template using tools-validator

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
TEMPLATE_PATH="${TEMPLATE_PATH:-}"
STRICT_MODE="${STRICT_MODE:-false}"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Validate template using tools-validator"
    echo ""
    echo "Options:"
    echo "  -p, --path PATH        Path to template (required)"
    echo "  -s, --strict           Enable strict mode"
    echo "  -h, --help             Show this help message"
    echo ""
    exit 1
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--path)
            TEMPLATE_PATH="$2"
            shift 2
            ;;
        -s|--strict)
            STRICT_MODE="true"
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
if [ -z "$TEMPLATE_PATH" ]; then
    echo -e "${RED}Error: Template path is required${NC}"
    usage
fi

TEMPLATE_PATH="$(cd "$(dirname "$TEMPLATE_PATH")" && pwd)/$(basename "$TEMPLATE_PATH")"

echo -e "${GREEN}Validating template${NC}"
echo "  Path: $TEMPLATE_PATH"
echo "  Strict mode: $STRICT_MODE"
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
echo -e "${GREEN}Executing template validation...${NC}"
cd "$PROJECT_ROOT"

CMD="$TOOLS_VALIDATOR -type template -path \"$TEMPLATE_PATH\""
[ "$STRICT_MODE" = "true" ] && CMD="$CMD -strict"

eval $CMD
VALIDATION_RESULT=$?

if [ $VALIDATION_RESULT -eq 0 ]; then
    echo -e "${GREEN}Template validation passed${NC}"
    exit 0
else
    echo -e "${RED}Error: Template validation failed${NC}"
    exit 1
fi
