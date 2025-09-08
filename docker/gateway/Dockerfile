FROM golang:1.25-alpine
RUN apk add --no-cache git tzdata bash busybox-extras \
 && go install github.com/air-verse/air@latest

WORKDIR /src/apps/gateway            
COPY apps/gateway/go.mod apps/gateway/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY apps/gateway/.air.toml ./
ENV APP_HOST=0.0.0.0 APP_PORT=8080
EXPOSE 8080

CMD ["air","-c",".air.toml"]