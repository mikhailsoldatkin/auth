FROM golang:1.22 AS builder

WORKDIR /auth

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

WORKDIR /auth

RUN mkdir -p cert

COPY --from=builder /auth/auth_server .

COPY .env .
COPY cert/service.key cert/service.key
COPY cert/service.pem cert/service.pem

CMD ["./auth_server"]
