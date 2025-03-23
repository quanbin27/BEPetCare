FROM golang:1.22.5 AS builder

WORKDIR /app/gateway
COPY ./gateway/go.mod ./gateway/go.sum ./
RUN go mod download
COPY ./gateway/*.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /gateway .
EXPOSE 8080
CMD ["./gateway"]
