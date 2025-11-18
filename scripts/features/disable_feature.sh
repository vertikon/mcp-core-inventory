#!/bin/bash
# Disable feature flag using yq

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
FEATURES_FILE="${PROJECT_ROOT}/config/features.yaml"
FEATURE_NAME="${FEATURE_NAME:-}"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Disable feature flag in features.yaml"
    echo ""
    echo "Options:"
    echo "  -f, --feature NAME     Feature name (required)"
    echo "  -h, --help             Show this help message"
    echo ""
    exit 1
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -f|--feature)
            FEATURE_NAME="$2"
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
if [ -z "$FEATURE_NAME" ]; then
    echo -e "${RED}Error: Feature name is required${NC}"
    usage
fi

# Check if yq is available
if ! command -v yq &> /dev/null; then
    echo -e "${RED}Error: yq is not installed. Install it from https://github.com/mikefarah/yq${NC}"
    exit 1
fi

# Check if features.yaml exists
if [ ! -f "$FEATURES_FILE" ]; then
    echo -e "${YELLOW}Warning: features.yaml not found${NC}"
    exit 1
fi

echo -e "${GREEN}Disabling feature: $FEATURE_NAME${NC}"

# Disable feature using yq
yq eval ".$FEATURE_NAME = false" -i "$FEATURES_FILE" || {
    echo -e "${RED}Error: Failed to disable feature${NC}"
    exit 1
}

echo -e "${GREEN}Feature '$FEATURE_NAME' disabled successfully${NC}"
