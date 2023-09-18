FROM golang:1.21.1-alpine3.17 as build-env

WORKDIR /go/src/app
ADD . /go/src/app

RUN go mod download && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /go/src/app/main/app /go/src/app/main

FROM alpine:3.18
COPY --from=build-env /go/src/app/main/app /app
COPY --from=build-env /go/src/app/configs/ /

CMD ["/app", "web", "-c", "/config.toml", "-m", "/model.conf", "--menu", "/menu.yaml"]
