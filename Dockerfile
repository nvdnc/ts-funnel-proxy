FROM golang:1.22-alpine3.19 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o ts-funnel-proxy .


FROM alpine:3.9

WORKDIR /app

COPY --from=builder /app/ts-funnel-proxy .

CMD ["./ts-funnel-proxy"]
