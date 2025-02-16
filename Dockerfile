FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM ubuntu:latest

COPY --from=builder . .

CMD ["/app/build/avito_app"]