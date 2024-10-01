FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o server .
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/migrate ./migrate
COPY start.sh .
COPY wait-for.sh .
COPY migrations ./migrations

EXPOSE 3001
CMD ["/app/server"]
ENTRYPOINT ["/app/start.sh"]
