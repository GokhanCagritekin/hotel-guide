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
RUN go build -o hotelguide ./cmd/hotelguide

# Final stage
FROM alpine:latest

# Uygulamayı çalıştırmak için gerekli çalışma dizinini ayarla
WORKDIR /root/

# Build aşamasında oluşturulan binary dosyasını kopyala
COPY --from=builder /app/hotelguide .

# Uygulamayı çalıştır
CMD ["./hotelguide"]
