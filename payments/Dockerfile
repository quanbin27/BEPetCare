FROM golang:1.23 AS builder

WORKDIR /app/payments
COPY ./payments/go.mod ./payments/go.sum ./
RUN go mod download
COPY ./payments/*.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /payments .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /payments .
EXPOSE 8086
CMD ["./payments"]
