FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Install specific swag version that works
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12

COPY . .

# Remove old docs and generate new ones
RUN rm -rf docs/docs.go docs/swagger.json docs/swagger.yaml
RUN swag init -g cmd/main.go

# Build with specific output name
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy binary
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]