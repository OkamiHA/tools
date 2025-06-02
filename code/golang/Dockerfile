# ---------- STAGE 1: Build Go binary ----------
FROM ubuntu:22.04 AS builder

# Cài Go và các công cụ cần thiết
RUN apt-get update && \
    apt-get install -y wget git build-essential curl ca-certificates

RUN wget https://go.dev/dl/go1.21.11.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.11.linux-amd64.tar.gz

ENV PATH="/usr/local/go/bin:$PATH"

# Đặt thư mục làm việc tại /app
WORKDIR /app

# Copy toàn bộ code vào /app
COPY . .

# Lấy module và build binary
RUN go mod tidy
RUN go build -o server main.go

# ---------- STAGE 2 ----------
FROM chromedp/headless-shell:latest

WORKDIR /app
ENV DEBIAN_FRONTEND=noninteractive
ENV LIBREOFFICE_VERSION=7.4.3.2
ENV LIBREOFFICE_DIR=LibreOffice_${LIBREOFFICE_VERSION}_Linux_x86-64_deb
# Cài công cụ và thư viện cần thiết
RUN apt-get update && apt-get install -y --no-install-recommends \
    wget ca-certificates tar gnupg \
    libxrender1 libxext6 libxinerama1 libsm6 libfontconfig1 libdbus-1-3 libgtk-3-0 libglu1-mesa libgdk-pixbuf2.0-0 libgl1-mesa-glx \
    fonts-dejavu fonts-liberation \
 && apt-get clean && rm -rf /var/lib/apt/lists/*

# Tải và cài LibreOffice từ archive
RUN wget https://downloadarchive.documentfoundation.org/libreoffice/old/${LIBREOFFICE_VERSION}/deb/x86_64/${LIBREOFFICE_DIR}.tar.gz && \
    tar -xvzf ${LIBREOFFICE_DIR}.tar.gz && \
    dpkg -i ${LIBREOFFICE_DIR}/DEBS/*.deb && \
    rm -rf ${LIBREOFFICE_DIR} ${LIBREOFFICE_DIR}.tar.gz

RUN ln -s /opt/libreoffice7.4/program/soffice /usr/local/bin/libreoffice

# Tạo thư mục output
RUN mkdir -p /app/output /app/converted

# Copy binary Go
COPY --from=builder /app/server /app/server

# Copy script start.sh và cấp quyền chạy
COPY start.sh /app/start.sh
COPY output /app/output
RUN chmod +x /app/start.sh


# Cổng Gin server
EXPOSE 8080

# Khởi động cả headless-shell và Gin
CMD ["/app/start.sh"]