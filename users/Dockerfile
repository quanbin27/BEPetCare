FROM golang:1.22.5 AS builder

WORKDIR /app/users
COPY ./users/go.mod ./users/go.sum ./
RUN go mod download
COPY ./users/*.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /users .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /users .
EXPOSE 8081
CMD ["./users"]
