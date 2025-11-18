#!/bin/bash
# Deploy using Docker using tools-deployer

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
PROJECT_NAME="${PROJECT_NAME:-}"
PROJECT_PATH="${PROJECT_PATH:-}"
IMAGE="${IMAGE:-}"
USE_COMPOSE="${USE_COMPOSE:-false}"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Deploy using Docker using tools-deployer"
    echo ""
    echo "Options:"
    echo "  -n, --name NAME        Project name (required)"
    echo "  -p, --path PATH        Project path (required)"
    echo "  -i, --image IMAGE      Container image (optional)"
    echo "  -c, --compose          Use docker-compose"
    echo "  -h, --help             Show this help message"
    echo ""
    exit 1
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--name)
            PROJECT_NAME="$2"
            shift 2
            ;;
        -p|--path)
            PROJECT_PATH="$2"
            shift 2
            ;;
        -i|--image)
            IMAGE="$2"
            shift 2
            ;;
        -c|--compose)
            USE_COMPOSE="true"
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
if [ -z "$PROJECT_NAME" ]; then
    echo -e "${RED}Error: Project name is required${NC}"
    usage
fi

if [ -z "$PROJECT_PATH" ]; then
    echo -e "${RED}Error: Project path is required${NC}"
    usage
fi

PROJECT_PATH="$(cd "$PROJECT_PATH" && pwd)"

echo -e "${GREEN}Deploying with Docker${NC}"
echo "  Project: $PROJECT_NAME"
echo "  Path: $PROJECT_PATH"
[ -n "$IMAGE" ] && echo "  Image: $IMAGE"
echo "  Use compose: $USE_COMPOSE"
echo ""

# Check if docker is available
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed${NC}"
    exit 1
fi

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

# Build tools-deployer if needed
TOOLS_DEPLOYER="${PROJECT_ROOT}/bin/tools-deployer"
if [ ! -f "$TOOLS_DEPLOYER" ]; then
    echo -e "${YELLOW}Building tools-deployer...${NC}"
    cd "$PROJECT_ROOT"
    mkdir -p bin
    go build -o "$TOOLS_DEPLOYER" ./cmd/tools-deployer || {
        echo -e "${RED}Error: Failed to build tools-deployer${NC}"
        exit 1
    }
fi

# Execute deployment
echo -e "${GREEN}Executing Docker deployment...${NC}"
cd "$PROJECT_ROOT"

CMD="$TOOLS_DEPLOYER -type docker -name \"$PROJECT_NAME\" -path \"$PROJECT_PATH\""
[ -n "$IMAGE" ] && CMD="$CMD -image \"$IMAGE\""

eval $CMD

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Docker deployment completed successfully${NC}"
else
    echo -e "${RED}Error: Docker deployment failed${NC}"
    exit 1
fi
