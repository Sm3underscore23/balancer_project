FROM  golang:1.24.2-alpine3.21 as builder

COPY . /app/source/
WORKDIR /app/source/

RUN go mod download
RUN go build -o ./bin/backend-test ./backend-test/

FROM alpine:3.21

WORKDIR /root/
COPY --from=builder /app/source/bin/backend-test .

ENV HOST_PORT=localhost:8090

CMD ["./backend-test"]
