FROM golang:1.24.2-alpine3.21 as builder

WORKDIR /app/backend-test
COPY ./backend-test/ ./


RUN go mod download
RUN go build -o /app/backend-test/bin/backend-test .

FROM alpine:3.21

WORKDIR /root/
COPY --from=builder /app/backend-test/bin/backend-test .

ENV HOST_PORT=localhost:8090

CMD ["./backend-test"]
