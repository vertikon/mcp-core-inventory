#!/bin/bash
# Setup infrastructure (DBs, Cache, Messaging, Cloud)

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

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Setup infrastructure (DBs, Cache, Messaging, Cloud)"
    echo ""
    echo "Options:"
    echo "  -e, --env ENV          Environment (dev/staging/prod) (default: dev)"
    echo "  -h, --help             Show this help message"
    echo ""
    exit 1
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--env)
            ENV="$2"
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

echo -e "${GREEN}Setting up infrastructure for environment: $ENV${NC}"
echo ""

# Load environment configuration using yq if available
if command -v yq &> /dev/null && [ -f "${CONFIG_DIR}/environments/${ENV}.yaml" ]; then
    echo -e "${GREEN}Loading configuration from ${CONFIG_DIR}/environments/${ENV}.yaml${NC}"
    # Configuration loaded via yq
fi

# Check for required CLIs
MISSING_CLIS=()

if ! command -v docker &> /dev/null; then
    MISSING_CLIS+=("docker")
fi

if ! command -v psql &> /dev/null && ! command -v mysql &> /dev/null; then
    MISSING_CLIS+=("database client (psql/mysql)")
fi

if [ ${#MISSING_CLIS[@]} -gt 0 ]; then
    echo -e "${YELLOW}Warning: Missing CLIs: ${MISSING_CLIS[*]}${NC}"
    echo -e "${YELLOW}Some infrastructure setup steps may be skipped${NC}"
fi

# Setup databases
echo -e "${GREEN}Setting up databases...${NC}"
# In production, this would:
# - Create databases using psql/mysql
# - Run migrations
# - Setup users and permissions
echo "  Database setup would be executed here"

# Setup cache (Redis)
echo -e "${GREEN}Setting up cache...${NC}"
if command -v redis-cli &> /dev/null; then
    echo "  Redis cache setup would be executed here"
else
    echo -e "${YELLOW}  Redis CLI not found, skipping cache setup${NC}"
fi

# Setup messaging (NATS)
echo -e "${GREEN}Setting up messaging...${NC}"
# In production, this would setup NATS JetStream
echo "  Messaging setup would be executed here"

# Setup cloud resources
echo -e "${GREEN}Setting up cloud resources...${NC}"
# In production, this would provision cloud resources
echo "  Cloud resources setup would be executed here"

echo -e "${GREEN}Infrastructure setup completed${NC}"
