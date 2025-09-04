# Hướng dẫn chạy với Docker

## Yêu cầu
- Cài đặt [Docker](https://docs.docker.com/get-docker/)  
- Cài đặt [Docker Compose](https://docs.docker.com/compose/install/)  

## Các lệnh Docker Compose

### Run Docker in WSL (without Docker Desktop)
```bash
sudo dockerd
```

### Build và chạy container
```bash
sudo docker compose up --build -d
```

### Dừng và xóa container
```bash
sudo docker compose down -v
```

### Tạm dừng container
```bash
sudo docker compose stop
```

### Khởi động lại container đã dừng
```bash
sudo docker compose start
```

---

## Ghi chú
