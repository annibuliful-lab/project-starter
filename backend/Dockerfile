# Stage 1: Build
FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# Stage 2: Run
FROM scratch

# Create folder structure to match godotenv.Load("../.env")
WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/.env /env/.env

ENTRYPOINT ["/app/main"]
