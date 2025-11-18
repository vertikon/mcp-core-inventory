#!/bin/bash
# Generate configuration files using tools-generator

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
CONFIG_NAME="${CONFIG_NAME:-}"
CONFIG_TYPE="${CONFIG_TYPE:-yaml}"
OUTPUT_PATH="${OUTPUT_PATH:-}"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Generate configuration files using tools-generator"
    echo ""
    echo "Options:"
    echo "  -n, --name NAME        Config name (required)"
    echo "  -t, --type TYPE        Config type (yaml/env/json/nats-schema/toml) (default: yaml)"
    echo "  -o, --output PATH      Output path (required)"
    echo "  -h, --help             Show this help message"
    echo ""
    exit 1
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--name)
            CONFIG_NAME="$2"
            shift 2
            ;;
        -t|--type)
            CONFIG_TYPE="$2"
            shift 2
            ;;
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

# Validate required parameters
if [ -z "$CONFIG_NAME" ]; then
    echo -e "${RED}Error: Config name is required${NC}"
    usage
fi

if [ -z "$OUTPUT_PATH" ]; then
    echo -e "${RED}Error: Output path is required${NC}"
    usage
fi

echo -e "${GREEN}Generating configuration file${NC}"
echo "  Name: $CONFIG_NAME"
echo "  Type: $CONFIG_TYPE"
echo "  Output: $OUTPUT_PATH"
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
echo -e "${GREEN}Executing config generation...${NC}"
cd "$PROJECT_ROOT"

CMD="$TOOLS_GENERATOR -type config -name \"$CONFIG_NAME\" -path \"$OUTPUT_PATH\" -stack \"$CONFIG_TYPE\""

eval $CMD

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Config generation completed successfully${NC}"
else
    echo -e "${RED}Error: Config generation failed${NC}"
    exit 1
fi
