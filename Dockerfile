# Dockerfile para mcp-core-inventory
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o core-inventory ./cmd/core-inventory

# Imagem final
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

# Copiar binário do builder
COPY --from=builder /app/core-inventory .

# Expor porta
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Executar aplicação
CMD ["./core-inventory"]

