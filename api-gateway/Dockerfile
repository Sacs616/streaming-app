# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy shared proto module first (for replace directive)
COPY shared/ /app/shared/

# Copy auth service (needed for proto types via replace directive)
COPY services/auth/ /app/services/auth/

# Copy api-gateway
COPY api-gateway/ /app/api-gateway/

WORKDIR /app/api-gateway

# Generate go.sum and download dependencies
RUN go mod tidy

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-gateway ./cmd/server

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/api-gateway/api-gateway .

EXPOSE 8080

CMD ["./api-gateway"]
