FROM  golang:1.24.2-alpine3.21 as builer

COPY . /app/source/
WORKDIR /app/source/

RUN go mod download
RUN go build -o ./bin/balancer cmd/main.go

FROM alpine:3.21

WORKDIR /root/
COPY --from=builer /app/source/bin/balancer .

CMD ["./balancer"]