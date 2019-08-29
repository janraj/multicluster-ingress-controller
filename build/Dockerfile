FROM golang:alpine AS builder
WORKDIR $GOPATH/src/multicluster-ingress-controller/
COPY pkg/   pkg
COPY vendor/  vendor
RUN go build -o /go/bin/multicluster-ingress-controller  ./pkg/main.go

FROM alpine
COPY --from=builder /go/bin/multicluster-ingress-controller  /go/bin/multicluster-ingress-controller
ENV PLATFORM Container
EXPOSE 8000
ENTRYPOINT ["/go/bin/multicluster-ingress-controller"]

