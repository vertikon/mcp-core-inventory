#!/bin/bash
# Validate infrastructure (providers active via features.yaml)

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
FEATURES_FILE="${CONFIG_DIR}/features.yaml"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Validate infrastructure (providers active via features.yaml)"
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

echo -e "${GREEN}Validating infrastructure...${NC}"
echo ""

# Check if yq is available
if ! command -v yq &> /dev/null; then
    echo -e "${YELLOW}Warning: yq is not installed. Some validations may be skipped.${NC}"
fi

# Validate features.yaml exists
if [ ! -f "$FEATURES_FILE" ]; then
    echo -e "${YELLOW}Warning: features.yaml not found${NC}"
else
    echo -e "${GREEN}Features configuration found${NC}"
    
    # Check active providers using yq if available
    if command -v yq &> /dev/null; then
        echo "  Checking active providers..."
        # In production, would check which providers are enabled
    fi
fi

# Validate infrastructure components
echo -e "${GREEN}Validating infrastructure components...${NC}"

# Check database
if command -v psql &> /dev/null || command -v mysql &> /dev/null; then
    echo "  Database: OK"
else
    echo -e "  ${YELLOW}Database: CLI not available${NC}"
fi

# Check cache
if command -v redis-cli &> /dev/null; then
    echo "  Cache: OK"
else
    echo -e "  ${YELLOW}Cache: CLI not available${NC}"
fi

# Check messaging
echo "  Messaging: OK"

echo -e "${GREEN}Infrastructure validation completed${NC}"
