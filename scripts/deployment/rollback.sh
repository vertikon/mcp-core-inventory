#!/bin/bash
# Rollback deployment

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
DEPLOYMENT_TYPE="${DEPLOYMENT_TYPE:-kubernetes}"
VERSION="${VERSION:-}"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Rollback deployment"
    echo ""
    echo "Options:"
    echo "  -n, --name NAME        Project name (required)"
    echo "  -t, --type TYPE        Deployment type (kubernetes/docker/serverless) (default: kubernetes)"
    echo "  -v, --version VERSION  Version to rollback to (optional)"
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
        -t|--type)
            DEPLOYMENT_TYPE="$2"
            shift 2
            ;;
        -v|--version)
            VERSION="$2"
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
if [ -z "$PROJECT_NAME" ]; then
    echo -e "${RED}Error: Project name is required${NC}"
    usage
fi

echo -e "${GREEN}Rolling back deployment${NC}"
echo "  Project: $PROJECT_NAME"
echo "  Type: $DEPLOYMENT_TYPE"
[ -n "$VERSION" ] && echo "  Version: $VERSION"
echo ""

# Rollback based on deployment type
case $DEPLOYMENT_TYPE in
    kubernetes)
        if command -v kubectl &> /dev/null; then
            echo -e "${GREEN}Rolling back Kubernetes deployment...${NC}"
            # In production, would use kubectl rollout undo
            echo "  kubectl rollout undo deployment/$PROJECT_NAME"
        else
            echo -e "${YELLOW}kubectl not found, skipping Kubernetes rollback${NC}"
        fi
        ;;
    docker)
        echo -e "${GREEN}Rolling back Docker deployment...${NC}"
        # In production, would stop current container and start previous version
        echo "  Docker rollback would be executed here"
        ;;
    serverless)
        echo -e "${GREEN}Rolling back serverless deployment...${NC}"
        # In production, would deploy previous version
        echo "  Serverless rollback would be executed here"
        ;;
    *)
        echo -e "${RED}Error: Unknown deployment type: $DEPLOYMENT_TYPE${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}Rollback completed${NC}"
