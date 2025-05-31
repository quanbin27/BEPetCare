FROM golang:1.23 AS builder

WORKDIR /app/products
COPY ./products/go.mod ./products/go.sum ./
RUN go mod download
COPY ./products/*.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /products .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /products .
EXPOSE 8082
CMD ["./products"]
