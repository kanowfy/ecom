FROM golang:1.19-alpine3.17 as builder
WORKDIR /app
COPY . .
RUN go build -o server .

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 3003
CMD ["/app/server"]
