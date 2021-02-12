FROM golang:alpine3.12

RUN mkdir -p /go/src/github.com/rtgnx/dcron
WORKDIR /go/src/github.com/rtgnx/dcron
COPY . .
RUN go build -o /usr/bin/dcron
CMD ["/usr/bin/dcron", "run", "--manifests", "/manifests/"]
