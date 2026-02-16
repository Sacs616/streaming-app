# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy shared proto module first (for replace directive)
COPY shared/ /app/shared/

# Copy auth service
COPY services/auth/ /app/services/auth/

WORKDIR /app/services/auth

# Generate go.sum and download dependencies
RUN go mod tidy

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-server ./cmd/server

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/auth/auth-server .

EXPOSE 50051

CMD ["./auth-server"]
