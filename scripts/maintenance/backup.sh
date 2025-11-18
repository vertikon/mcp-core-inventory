#!/bin/bash
# Backup data

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
BACKUP_DIR="${BACKUP_DIR:-${PROJECT_ROOT}/backups}"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Backup data"
    echo ""
    echo "Options:"
    echo "  -e, --env ENV          Environment (dev/staging/prod) (default: dev)"
    echo "  -d, --dir DIR          Backup directory (default: ./backups)"
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
        -d|--dir)
            BACKUP_DIR="$2"
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

echo -e "${GREEN}Creating backup for environment: $ENV${NC}"
echo "  Backup directory: $BACKUP_DIR"
echo ""

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Backup databases
echo -e "${GREEN}Backing up databases...${NC}"
# In production, would backup actual databases
echo "  Database backup would be executed here"

# Backup configuration
echo -e "${GREEN}Backing up configuration...${NC}"
if [ -d "$CONFIG_DIR" ]; then
    BACKUP_CONFIG_DIR="${BACKUP_DIR}/config-$(date +%Y%m%d-%H%M%S)"
    cp -r "$CONFIG_DIR" "$BACKUP_CONFIG_DIR"
    echo "  Configuration backed up to: $BACKUP_CONFIG_DIR"
fi

# Backup state
echo -e "${GREEN}Backing up state...${NC}"
# In production, would backup state store
echo "  State backup would be executed here"

echo -e "${GREEN}Backup completed${NC}"
