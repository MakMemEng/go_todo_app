# デプロイ用コンテナに含めるバイナリを作成
FROM golang:1.21.0-bullseye as deploy-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -trimpath -ldflags "-w -s" -o app

# --------------------------------------------

# デプロイ用コンテナ
# docker build -t makmemeng/ gotodo:${DOCKER_TAG} --target deploy ./
FROM debian:bullseye-slim as deploy

RUN apt-get update

COPY --from=deploy-builder /app/app .

CMD [ "./app" ]

# --------------------------------------------


# ローカル開発環境で利用するホットリロード環境
FROM golang:1.21.0 as dev
WORKDIR /app
RUN go install github.com/cosmtrek/air@latest
CMD [ "air" ]