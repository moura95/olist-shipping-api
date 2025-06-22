FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Try building with existing docs first
RUN go build cmd/main.go || echo "Build failed, will try without docs"

FROM alpine
COPY --from=builder /app/main /usr/local/bin/main
COPY --from=builder /app/.env /usr/local/bin/.env
EXPOSE 8080
WORKDIR /usr/local/bin/
ENTRYPOINT ["main"]