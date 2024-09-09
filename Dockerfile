FROM golang:1.21.6-bookworm

WORKDIR /app
COPY go.mod main.go /app/
COPY tpl/ /app/tpl/

RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/outbox-relay

ENTRYPOINT ["/app/outbox-relay"]
CMD ["help"]

