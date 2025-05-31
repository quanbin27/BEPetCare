FROM golang:1.23 AS builder

WORKDIR /app/gateway
COPY ./gateway/go.mod ./gateway/go.sum ./
COPY ./gateway/docs ./docs
RUN go mod download
COPY ./gateway/*.go ./
COPY ./gateway/handlers ./handlers
RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /gateway .
EXPOSE 8080
CMD ["./gateway"]
