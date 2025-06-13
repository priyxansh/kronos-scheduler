# Build stage
FROM golang:1.24.2 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler ./cmd/scheduler/main.go

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/scheduler .
COPY .env .env
EXPOSE 3000
ENTRYPOINT ["./scheduler"]
