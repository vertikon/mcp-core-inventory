#!/bin/bash
# Migrate data versions and schemas between environments

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
SOURCE_ENV="${SOURCE_ENV:-dev}"
TARGET_ENV="${TARGET_ENV:-staging}"
DATASET_ID="${DATASET_ID:-}"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Migrate data versions and schemas between environments"
    echo ""
    echo "Options:"
    echo "  -s, --source-env SOURCE    Source environment (default: dev)"
    echo "  -t, --target-env TARGET    Target environment (default: staging)"
    echo "  -d, --dataset-id ID        Dataset ID (required)"
    echo "  -h, --help                 Show this help message"
    echo ""
    exit 1
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -s|--source-env)
            SOURCE_ENV="$2"
            shift 2
            ;;
        -t|--target-env)
            TARGET_ENV="$2"
            shift 2
            ;;
        -d|--dataset-id)
            DATASET_ID="$2"
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
if [ -z "$DATASET_ID" ]; then
    echo -e "${RED}Error: Dataset ID is required${NC}"
    usage
fi

echo -e "${GREEN}Starting data migration${NC}"
echo "  Source: $SOURCE_ENV"
echo "  Target: $TARGET_ENV"
echo "  Dataset ID: $DATASET_ID"
echo ""

# Execute migration
echo -e "${GREEN}Executing migration...${NC}"
cd "$PROJECT_ROOT"

# Load configuration for source and target environments
CONFIG_DIR="${PROJECT_ROOT}/config"
SOURCE_CONFIG="${CONFIG_DIR}/environments/${SOURCE_ENV}.yaml"
TARGET_CONFIG="${CONFIG_DIR}/environments/${TARGET_ENV}.yaml"

# In production, this would call the migration engine from internal/versioning/data
# For now, validate configuration and prepare migration
if [ ! -f "$SOURCE_CONFIG" ]; then
    echo -e "${YELLOW}Warning: Source environment config not found: $SOURCE_CONFIG${NC}"
fi

if [ ! -f "$TARGET_CONFIG" ]; then
    echo -e "${YELLOW}Warning: Target environment config not found: $TARGET_CONFIG${NC}"
fi

# Migration would be executed via Go migration engine
# This requires creating cmd/migration-data/main.go that uses internal/versioning/data
echo "  - Source environment: $SOURCE_ENV"
echo "  - Target environment: $TARGET_ENV"
echo "  - Dataset ID: $DATASET_ID"
echo "  - Migration engine: internal/versioning/data"

# Note: Full implementation requires cmd/migration-data CLI tool
echo -e "${GREEN}Migration preparation completed${NC}"
echo -e "${YELLOW}Note: Full migration execution requires cmd/migration-data CLI tool${NC}"
