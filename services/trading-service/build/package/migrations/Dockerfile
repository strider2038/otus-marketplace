FROM golang:1-alpine AS build-env

RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@master

FROM scratch

ENV GOOSE_DRIVER="postgres" \
    GOOSE_DBSTRING=""

COPY --from=build-env /go/bin/goose /go/bin/goose
COPY ./scripts/migrations/ /migrations/

WORKDIR /migrations

CMD ["/go/bin/goose", "up", "-v"]
