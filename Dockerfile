FROM golang:1.17.3-alpine3.14 as build-env

ARG VERSION

WORKDIR /go/src/app
ADD . /go/src/app

RUN go mod download && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.VERSION=${VERSION}" -o /go/src/app/main/app /go/src/app/main

FROM alpine:3.14
COPY --from=build-env /go/src/app/main/app /app

CMD ["/app", "web", "-c", "/config.toml", "-m", "/model.conf", "--menu", "/menu.yaml"]
