FROM golang:1.10 as builder

ARG LDFLAGS=""

ENV GO_PATH=/go/src/github.com/derek-elliott/url-shortener
ENV CGO_ENABLED=0
ENV GOGC=off
ENV GOOS=linux

COPY . $GO_PATH
WORKDIR $GO_PATH

RUN go get -u github.com/golang/dep/cmd/dep && \
    dep ensure && \
    go build -a -installsuffix nocgo -o snip -ldflags "$LDFLAGS" .

FROM scratch
COPY --from=builder /go/src/github.com/derek-elliott/url-shortener/snip /
CMD ["/snip"]
