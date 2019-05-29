FROM golang:alpine AS builder
WORKDIR $GOPATH/src/k8s-control-plane/
COPY cmd/   cmd
COPY handler/	handler
COPY vendor/  vendor
RUN go build -o /go/bin/k8s-control-plane  ./cmd/k8s-control-plane/main.go

FROM alpine
COPY --from=builder /go/bin/k8s-control-plane  /go/bin/k8s-control-plane
EXPOSE 8080
ENTRYPOINT ["/go/bin/k8s-control-plane"]

