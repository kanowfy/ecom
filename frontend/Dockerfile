FROM golang:1.22-alpine3.20 as builder
WORKDIR /app
COPY . .
RUN go build -o server .

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/ui ./ui
EXPOSE 4000
CMD ["/app/server"]
