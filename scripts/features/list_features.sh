#!/bin/bash
# List feature flags using yq

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

echo -e "${GREEN}Feature Flags:${NC}"
echo ""

# List features using yq
yq eval 'to_entries | .[] | "\(.key): \(.value)"' "$FEATURES_FILE" || {
    echo -e "${YELLOW}No features found${NC}"
}
