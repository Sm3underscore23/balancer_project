FROM  golang:1.24.2-alpine3.21 as builder

COPY . /app/source/
WORKDIR /app/source/

RUN go mod download
RUN go build -o ./bin/balancer ./cmd/main.go

FROM alpine:3.21

WORKDIR /root/
COPY --from=builder /app/source/bin/balancer .

CMD ["./balancer", "config-path=./config/config.yaml"]