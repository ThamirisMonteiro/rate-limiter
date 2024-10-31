FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o rate-limiter ./cmd/ratelimiter/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/rate-limiter .

EXPOSE 8080

RUN ls -la .

CMD ["./rate-limiter"]
