FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/backend ./cmd/main.go

FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/backend /app/backend
COPY --from=builder /app/migrations /app/migrations
COPY --from=builder /app/docs /app/docs
COPY --from=builder /app/.env /app/.env
EXPOSE 3000
CMD ["/app/backend"]