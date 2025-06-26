FROM golang:1.24.2-alpine3.21 AS builder

WORKDIR /app/source/
COPY . .

RUN go mod download
RUN mkdir -p ./bin
RUN go build -o ./bin/balancer ./cmd/main.go

FROM alpine:3.21

WORKDIR /root/
COPY --from=builder /app/source/bin/balancer ./

CMD ["./balancer", "-config-path=config.yaml"]
