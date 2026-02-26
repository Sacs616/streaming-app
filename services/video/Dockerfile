FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy shared proto module first (for replace directive)
COPY shared/ /app/shared/

# Copy auth service
COPY services/video/ /app/services/video/

WORKDIR /app/services/video

# Generate go.sum and download dependencies
RUN go mod tidy

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o video-service ./cmd/server

# Runtime stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/services/video/video-service .

EXPOSE 50052

CMD ["./video-service"]