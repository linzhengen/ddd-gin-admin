FROM golang:1.16 as build-env

WORKDIR /go/src/app
ADD . /go/src/app

RUN go mod download
RUN make build

FROM gcr.io/distroless/base
COPY --from=build-env /go/src/app/main/app /app
CMD ["/app"]
