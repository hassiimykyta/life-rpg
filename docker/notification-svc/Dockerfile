# syntax=docker/dockerfile:1.7
FROM golang:1.25-alpine

# утилиты + air
RUN apk add --no-cache git tzdata bash busybox-extras \
 && go install github.com/air-verse/air@latest

# работаем внутри модуля auth-svc
WORKDIR /src/apps/auth-svc

# прогреем зависимости по модулю (ускоряет старт Air)
COPY apps/auth-svc/go.mod apps/auth-svc/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

# конфиг Air (сам код смонтируем томом)
COPY apps/auth-svc/.air.toml ./

# сетап
ENV GRPC_PORT=8082
EXPOSE 8082

# стартуем hot-reload
CMD ["air","-c",".air.toml"]