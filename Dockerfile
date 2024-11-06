FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download


RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

FROM alpine

WORKDIR /app

COPY --from=builder /app/app /app/app
COPY --from=builder /app/config.yaml /app/config.yaml
COPY --from=builder /app/migrations /app/migrations

CMD ["/app/app"]