#!/bin/bash
# Generate documentation (OpenAPI/AsyncAPI/Schemas)

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Generate documentation (OpenAPI/AsyncAPI/Schemas)"
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

echo -e "${GREEN}Generating documentation...${NC}"
echo ""

# Generate OpenAPI
echo -e "${GREEN}Generating OpenAPI...${NC}"
"${SCRIPT_DIR}/generate_openapi.sh" || echo -e "${YELLOW}OpenAPI generation skipped${NC}"

# Generate AsyncAPI
echo -e "${GREEN}Generating AsyncAPI...${NC}"
"${SCRIPT_DIR}/generate_asyncapi.sh" || echo -e "${YELLOW}AsyncAPI generation skipped${NC}"

echo -e "${GREEN}Documentation generation completed${NC}"
