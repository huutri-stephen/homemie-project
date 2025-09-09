# Dockerfile

# Stage 1: build
FROM golang:1.23-alpine AS builder
# Dùng image golang phiên bản 1.23 trên nền alpine (nhẹ, tối giản)
# Đây là môi trường build chính, chứa compiler Go và các tool cần thiết
WORKDIR /app
# Thiết lập thư mục làm việc mặc định trong container là /app
# Tất cả các lệnh tiếp theo sẽ được chạy trong thư mục này
COPY go.mod go.sum ./
# Copy 2 file go.mod và go.sum từ local vào /app trong container
COPY vendor vendor
# Copy vendored dependencies
COPY . .
# Copy toàn bộ nội dung thư mục hiện tại vào container tại /app
# Mang toàn bộ mã nguồn lên container để build và chạy.
RUN go build -mod=vendor -o homemie cmd/main.go
# Biên dịch file cmd/main.go thành file nhị phân tên homemie
# Tạo executable để container chạy khi khởi động

# Stage 2: run
FROM alpine:3.20
WORKDIR /root/
COPY --from=builder /app/homemie .
COPY --from=builder /app/.env .
EXPOSE 8080
# Khai báo port 8080 để container expose ra ngoài
CMD ["./homemie"]
# Khi container chạy, nó sẽ thực thi binary ./homemie đã build ở trên