FROM golang:1.20.1-alpine3.17 as build-env

ARG VERSION

WORKDIR /go/src/app
ADD . /go/src/app

RUN go mod download && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X main.VERSION=${VERSION}" -o /go/src/app/main/app /go/src/app/main

FROM alpine:3.17
COPY --from=build-env /go/src/app/main/app /app

CMD ["/app", "web", "-c", "/config.toml", "-m", "/model.conf", "--menu", "/menu.yaml"]
