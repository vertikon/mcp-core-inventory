#!/bin/bash
# Deploy to Kubernetes using tools-deployer

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

# Default values
PROJECT_NAME="${PROJECT_NAME:-}"
PROJECT_PATH="${PROJECT_PATH:-}"
NAMESPACE="${NAMESPACE:-default}"
IMAGE="${IMAGE:-}"
MANIFESTS_PATH="${MANIFESTS_PATH:-}"
REPLICAS="${REPLICAS:-1}"
KUBECONFIG="${KUBECONFIG:-}"

usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Deploy to Kubernetes using tools-deployer"
    echo ""
    echo "Options:"
    echo "  -n, --name NAME        Project name (required)"
    echo "  -p, --path PATH        Project path (required)"
    echo "  --namespace NS         Kubernetes namespace (default: default)"
    echo "  -i, --image IMAGE      Container image (required if no manifests)"
    echo "  -m, --manifests PATH   Path to Kubernetes manifests"
    echo "  -r, --replicas N       Number of replicas (default: 1)"
    echo "  --kubeconfig PATH      Path to kubeconfig file"
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
        --namespace)
            NAMESPACE="$2"
            shift 2
            ;;
        -i|--image)
            IMAGE="$2"
            shift 2
            ;;
        -m|--manifests)
            MANIFESTS_PATH="$2"
            shift 2
            ;;
        -r|--replicas)
            REPLICAS="$2"
            shift 2
            ;;
        --kubeconfig)
            KUBECONFIG="$2"
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

if [ -z "$PROJECT_PATH" ]; then
    echo -e "${RED}Error: Project path is required${NC}"
    usage
fi

if [ -z "$IMAGE" ] && [ -z "$MANIFESTS_PATH" ]; then
    echo -e "${RED}Error: Image or manifests path is required${NC}"
    usage
fi

# Resolve absolute path
PROJECT_PATH="$(cd "$PROJECT_PATH" && pwd)"

echo -e "${GREEN}Deploying to Kubernetes${NC}"
echo "  Project: $PROJECT_NAME"
echo "  Path: $PROJECT_PATH"
echo "  Namespace: $NAMESPACE"
[ -n "$IMAGE" ] && echo "  Image: $IMAGE"
[ -n "$MANIFESTS_PATH" ] && echo "  Manifests: $MANIFESTS_PATH"
echo "  Replicas: $REPLICAS"
echo ""

# Check if kubectl is available (optional, but recommended)
if ! command -v kubectl &> /dev/null; then
    echo -e "${YELLOW}Warning: kubectl is not installed. Deployment may fail.${NC}"
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
echo -e "${GREEN}Executing Kubernetes deployment...${NC}"
cd "$PROJECT_ROOT"

# Build command
CMD="$TOOLS_DEPLOYER -type kubernetes -name \"$PROJECT_NAME\" -path \"$PROJECT_PATH\" -namespace \"$NAMESPACE\" -replicas $REPLICAS"
[ -n "$IMAGE" ] && CMD="$CMD -image \"$IMAGE\""
[ -n "$MANIFESTS_PATH" ] && CMD="$CMD -manifests \"$MANIFESTS_PATH\""
[ -n "$KUBECONFIG" ] && CMD="$CMD -kubeconfig \"$KUBECONFIG\""

# Execute
eval $CMD

if [ $? -eq 0 ]; then
    echo -e "${GREEN}Kubernetes deployment completed successfully${NC}"
else
    echo -e "${RED}Error: Kubernetes deployment failed${NC}"
    exit 1
fi
