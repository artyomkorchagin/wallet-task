FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/wallet-service ./cmd/main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/backend /app/backend
COPY  ./migrations /app/migrations
COPY ./docs /app/docs
COPY ./config.env /app/config.env
COPY ./postgresql.conf ./db/postgresql.conf
EXPOSE 3000
CMD ["/app/wallet-service"]