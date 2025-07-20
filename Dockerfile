# Dockerfile for building a Go application
FROM golang:1.24.5-alpine AS builder

# 作業ディレクトリの設定
WORKDIR /app

# 依存関係をダウンロード
COPY go.mod ./
RUN go mod download

# アプリケーションのソースコードをコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o main ./cmd/server

# 実行環境のイメージ
FROM alpine:latest

# 作業ディレクトリの設定
WORKDIR /app

# ビルドステージからコンパイル済みのバイナリをコピー
COPY --from=builder /app/main .

# アプリケーションがリッスンするポートを公開
EXPOSE 8080

# コンテナ起動時にGoアプリケーションを実行
CMD ["./main"]
