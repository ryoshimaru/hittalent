FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app ./cmd/app
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.21.1


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app /app/app
COPY --from=builder /go/bin/goose /app/goose
COPY migrations /app/migrations

EXPOSE 8080

CMD ["/app/app"]