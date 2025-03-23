FROM golang:1.22.5 AS builder

WORKDIR /app/orders
COPY ./orders/go.mod ./orders/go.sum ./
RUN go mod download
COPY ./orders/*.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /orders .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /orders .
EXPOSE 8083
CMD ["./orders"]
