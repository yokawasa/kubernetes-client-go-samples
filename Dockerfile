## multi-stage build
# Stage - Binary build
FROM golang:1.14.2-buster as builder
RUN apt update && apt install -y --no-install-recommends git tzdata
WORKDIR /go/src/github.com/yokawasa/kubernetes-client-go-sample
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
COPY . .
RUN make

# Stage - Runtime
FROM debian:10.3-slim as executor
RUN apt-get update && apt-get install -y python python-pip \
 && rm -rf /var/lib/apt/lists/* \
 && pip install awscli
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /
COPY --from=builder /go/src/github.com/yokawasa/kubernetes-client-go-sample/dist/svclist /svclist
COPY --from=builder /go/src/github.com/yokawasa/kubernetes-client-go-sample/dist/nodelist /nodelist
COPY --from=builder /go/src/github.com/yokawasa/kubernetes-client-go-sample/dist/podlist /podlist
COPY --from=builder /go/src/github.com/yokawasa/kubernetes-client-go-sample/dist/podlist-in-svc /podlist-in-svc
COPY --from=builder /go/src/github.com/yokawasa/kubernetes-client-go-sample/dist/informer /informer
COPY --from=builder /go/src/github.com/yokawasa/kubernetes-client-go-sample/dist/istio-get-vs /istio-get-vs
COPY --from=builder /go/src/github.com/yokawasa/kubernetes-client-go-sample/dist/istio-update-vs /istio-update-vs
