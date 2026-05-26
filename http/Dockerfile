FROM golang:1.26.3-alpine3.22 AS builder

WORKDIR /app

COPY . .

RUN apk add --no-cache ca-certificates

RUN go build -o http ./cmd

FROM scratch

WORKDIR /app

COPY --from=builder /app/http .
COPY --from=builder /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /etc/ssl/cert.pem /etc/ssl/cert.pem

EXPOSE 8099

CMD ["./http"]