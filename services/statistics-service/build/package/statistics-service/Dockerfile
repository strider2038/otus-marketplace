FROM golang:1-alpine AS build-env

ARG APP_VERSION=""

COPY . /project

RUN set -xe \
    && cd /project/cmd/statistics-service \
    && go mod download \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s -X main.version=${APP_VERSION}" -o statistics-service

FROM scratch

ENV DATABASE_URL=""

COPY --from=build-env /project/cmd/statistics-service/statistics-service /app/

EXPOSE 8000

ENTRYPOINT ["/app/statistics-service"]
