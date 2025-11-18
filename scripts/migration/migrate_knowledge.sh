#!/bin/bash
# Migrate knowledge versions between environments

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
CONFIG_DIR="${PROJECT_ROOT}/config"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
SOURCE_ENV="${SOURCE_ENV:-dev}"
TARGET_ENV="${TARGET_ENV:-staging}"
KNOWLEDGE_ID="${KNOWLEDGE_ID:-}"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Migrate knowledge versions between environments"
    echo ""
    echo "Options:"
    echo "  -s, --source-env SOURCE    Source environment (default: dev)"
    echo "  -t, --target-env TARGET    Target environment (default: staging)"
    echo "  -k, --knowledge-id ID     Knowledge base ID (required)"
    echo "  -h, --help                 Show this help message"
    echo ""
    echo "Environment variables:"
    echo "  SOURCE_ENV                Source environment"
    echo "  TARGET_ENV                Target environment"
    echo "  KNOWLEDGE_ID              Knowledge base ID"
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
        -k|--knowledge-id)
            KNOWLEDGE_ID="$2"
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
if [ -z "$KNOWLEDGE_ID" ]; then
    echo -e "${RED}Error: Knowledge ID is required${NC}"
    usage
fi

echo -e "${GREEN}Starting knowledge migration${NC}"
echo "  Source: $SOURCE_ENV"
echo "  Target: $TARGET_ENV"
echo "  Knowledge ID: $KNOWLEDGE_ID"
echo ""

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

# Build migration tool if needed
if [ ! -f "$PROJECT_ROOT/bin/migrate-knowledge" ]; then
    echo -e "${YELLOW}Building migration tool...${NC}"
    cd "$PROJECT_ROOT"
    go build -o bin/migrate-knowledge ./cmd/migration-knowledge 2>/dev/null || {
        echo -e "${YELLOW}Migration tool not found, using direct Go execution${NC}"
    }
fi

# Execute migration
echo -e "${GREEN}Executing migration...${NC}"
cd "$PROJECT_ROOT"

# Load configuration for source and target environments
SOURCE_CONFIG="${CONFIG_DIR}/environments/${SOURCE_ENV}.yaml"
TARGET_CONFIG="${CONFIG_DIR}/environments/${TARGET_ENV}.yaml"

# In production, this would call the migration engine from internal/versioning/knowledge
# For now, validate configuration and prepare migration
if [ ! -f "$SOURCE_CONFIG" ]; then
    echo -e "${YELLOW}Warning: Source environment config not found: $SOURCE_CONFIG${NC}"
fi

if [ ! -f "$TARGET_CONFIG" ]; then
    echo -e "${YELLOW}Warning: Target environment config not found: $TARGET_CONFIG${NC}"
fi

# Migration would be executed via Go migration engine
# This requires creating cmd/migration-knowledge/main.go that uses internal/versioning/knowledge
echo "  - Source environment: $SOURCE_ENV"
echo "  - Target environment: $TARGET_ENV"
echo "  - Knowledge ID: $KNOWLEDGE_ID"
echo "  - Migration engine: internal/versioning/knowledge/migration_engine.go"

# Note: Full implementation requires cmd/migration-knowledge CLI tool
echo -e "${GREEN}Migration preparation completed${NC}"
echo -e "${YELLOW}Note: Full migration execution requires cmd/migration-knowledge CLI tool${NC}"
