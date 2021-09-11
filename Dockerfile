FROM golang:latest as build

WORKDIR /root
ADD . /go/src/app

WORKDIR /go/src/app
ENV CGO_ENABLED=0
RUN go build -o /go/src/app/dynu-dns-updater /go/src/app/main.go


FROM alpine:latest

COPY --from=build /go/src/app/dynu-dns-updater /usr/local/bin/dynu-dns-updater

CMD /usr/local/bin/dynu-dns-updater