# Stage 1: Builder
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o comment-tree ./cmd/app/main.go

# Stage 2: Runner
FROM alpine:latest

WORKDIR /root/

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/comment-tree .

COPY --from=builder /app/web ./web

COPY --from=builder /app/configs ./configs

COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./comment-tree"]