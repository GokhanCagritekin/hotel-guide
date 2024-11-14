# Build stage
FROM golang:1.20-alpine AS builder

# Proje dosyalarını içerisine kopyalayacağımız çalışma dizinini ayarla
WORKDIR /app

# Go mod dosyalarını yükleyin ve bağımlılıkları indirin
COPY go.mod go.sum ./
RUN go mod download

# Proje dosyalarının tamamını kopyala
COPY . .

# Uygulamayı build et (binary dosyasını oluştur)
RUN go build -o hotelservice ./cmd/hotel-service/main.go

# Final stage
FROM alpine:latest

# Uygulamayı çalıştırmak için gerekli çalışma dizinini ayarla
WORKDIR /root/

# Build aşamasında oluşturulan binary dosyasını kopyala
COPY --from=builder /app/hotelservice .

# Uygulamayı çalıştır
CMD ["./hotelservice"]
