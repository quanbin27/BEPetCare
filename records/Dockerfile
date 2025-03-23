FROM golang:1.22.5 AS builder

WORKDIR /app/records
COPY ./records/go.mod ./records/go.sum ./
RUN go mod download
COPY ./records/*.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /records .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /records .
EXPOSE 8085
CMD ["./records"]
