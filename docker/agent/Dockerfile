FROM golang:alpine AS builder

WORKDIR /build

COPY ./calculator .

RUN go build -o main ./cmd/agent/

FROM alpine

WORKDIR /app

COPY --from=builder /build/main .

ENTRYPOINT ["./main"]
