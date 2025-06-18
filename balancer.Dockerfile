FROM  golang:1.24.2-alpine3.21 as builder

WORKDIR /app/source/
COPY ./cmd ./config ./internal ./go.mod ./go.sum ./

RUN go mod download
RUN go build -o ./bin/balancer ./cmd/main.go

FROM alpine:3.21

WORKDIR /root/
COPY --from=builder /app/source/bin/balancer .

CMD ["./balancer", "-config-path=./config/config.yaml"]
