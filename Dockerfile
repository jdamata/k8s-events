FROM golang AS build-env

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/github.com/jdamata/k8s-events
ADD . /go/src/github.com/jdamata/k8s-events
RUN go build -a -tags netgo -ldflags '-w' -o /bin/k8s-events

FROM alpine
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /bin/k8s-events /k8s-events
ENTRYPOINT ["/k8s-events"]