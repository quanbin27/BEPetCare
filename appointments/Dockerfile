FROM golang:1.22.5 AS builder

# Đặt thư mục làm việc
WORKDIR /app/appointments

# Copy go.mod và go.sum của appointment service
COPY ./appointments/go.mod ./appointments/go.sum ./
COPY .env .env
# Download dependencies
RUN go mod download

# Copy source code của appointment service
COPY ./appointments/*.go ./

# Build ứng dụng
RUN CGO_ENABLED=0 GOOS=linux go build -o /appointments .

# Sử dụng alpine làm base image cho runtime stage
FROM alpine:latest

# Cài đặt SSL certificates
RUN apk --no-cache add ca-certificates

# Đặt thư mục làm việc
WORKDIR /root/

# Copy binary từ builder stage
COPY --from=builder /appointments .

# Expose port mà service sẽ lắng nghe
EXPOSE 8084

# Command để chạy ứng dụng
CMD ["./appointments"]