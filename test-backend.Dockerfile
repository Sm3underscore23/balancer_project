FROM  golang:1.24.2-alpine3.21 as builder

COPY . /app/source/
WORKDIR /app/source/

RUN go mod download
RUN go build -o ./bin/test-backend ./test-backend/

FROM alpine:3.21

WORKDIR /root/
COPY --from=builder /app/source/bin/test-backend .

ENV HOST_PORT=localhost:8090

CMD ["./test-backend"]