FROM golang:1.22.5 AS builder

WORKDIR /app/notification
COPY ./notification/go.mod ./notification/go.sum ./
RUN go mod download
COPY ./notification/*.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /notification .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /notification .
EXPOSE 8077
CMD ["./notification"]
