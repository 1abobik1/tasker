FROM golang:1.21.7-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o tasker ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/tasker .
COPY --from=builder /app/.env .

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

CMD ["./tasker", "--config=.env"]