FROM golang:1.21.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

FROM ubuntu:latest

COPY --from=builder . .

#CMD sleep infinity
CMD ["/app/build/avito_app"]