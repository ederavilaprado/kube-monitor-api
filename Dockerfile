FROM golang:1.6.3-wheezy

RUN mkdir -p /go/src/github.com/ederavilaprado/kube-monitor
WORKDIR /go/src/github.com/ederavilaprado/kube-monitor
COPY . /go/src/github.com/ederavilaprado/kube-monitor

RUN go get github.com/tools/godep
RUN godep go install
CMD ["kube-monitor"]
EXPOSE 5000
