# Dockerfile para mcp-core-inventory - Otimizado para Produção
FROM golang:1.23-alpine AS builder

# Instalar dependências necessárias
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copiar go mod files
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copiar código fonte
COPY . .

# Build da aplicação com otimizações
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o mcp-core-inventory ./cmd

# Imagem final minimalista
FROM alpine:latest

# Instalar dependências de runtime
RUN apk --no-cache add ca-certificates tzdata curl && \
    adduser -D -s /bin/sh appuser

WORKDIR /app

# Copiar binário do builder
COPY --from=builder /app/mcp-core-inventory .
COPY --from=builder /app/config ./config

# Mudar para usuário não-root
USER appuser

# Expor porta
EXPOSE 8080

# Health check aprimorado
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Executar aplicação
CMD ["./mcp-core-inventory"]

