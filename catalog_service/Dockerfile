FROM golang:1.19-alpine3.16 as builder
WORKDIR /app
COPY . .
RUN go build -o server main.go

FROM alpine3.16
WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 3001
CMD ["/app/server"]

