FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o server .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 3004
CMD ["/app/server"]
