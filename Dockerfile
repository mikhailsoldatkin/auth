FROM golang:1.22 as builder

WORKDIR /app

RUN apt-get update &&  \
    apt-get upgrade -y &&  \
    apt-get clean && rm -rf /var/lib/apt/lists/*

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o auth_server ./cmd/grpc_server/main.go

FROM alpine:latest

RUN apk update && \
    apk upgrade && \
    rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=builder /app/auth_server .
COPY .env .

CMD ["./auth_server"]
