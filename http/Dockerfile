FROM golang:1.26.3-alpine3.22 AS builder

WORKDIR /app

COPY . .

RUN go build -o http ./cmd

FROM gcr.io/distroless/static-debian13

WORKDIR /app

COPY --from=builder /app/http .

EXPOSE 8099

CMD ["./http"]