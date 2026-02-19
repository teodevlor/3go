# ========= Build stage =========
FROM golang:1.22-alpine AS builder
ENV GOTOOLCHAIN=auto
WORKDIR /gogogo-api

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o app ./cmd/server/main.go


# ========= Runtime stage =========
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# copy binary
COPY --from=builder /gogogo-api/app .

# copy config
COPY --from=builder /gogogo-api/config ./config

EXPOSE 8080

CMD ["./app"]
