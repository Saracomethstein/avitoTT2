FROM golang:1.22 AS builder

WORKDIR /app

COPY .env ./
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/build/avito_app /app/avito_app
COPY --from=builder /app/.env /app/.env

CMD ["/app/avito_app"]