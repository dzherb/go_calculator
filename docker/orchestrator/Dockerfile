FROM golang:alpine AS builder

WORKDIR /build

COPY ./calculator_services .

RUN go build -o main ./cmd/orchestrator/

FROM alpine

WORKDIR /app

COPY --from=builder /build/main .

EXPOSE 8080

ENTRYPOINT ["./main"]
