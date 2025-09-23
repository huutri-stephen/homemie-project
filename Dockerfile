# Stage 1: build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o homemie ./cmd/main.go

# Stage 2: run
FROM alpine:3.20
WORKDIR /root/
COPY --from=builder /app/homemie .
EXPOSE 8080
CMD ["./homemie"]
