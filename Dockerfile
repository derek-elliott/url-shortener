FROM golang:1.8.5-jessie as builder

ARG LDFLAGS=""

ENV GO_PATH=/go/src/github.com/derek-elliott/url-shortener

COPY . $GO_PATH

WORKDIR $GO_PATH

RUN go get -u github.com/golang/dep/cmd/dep && \
    dep ensure && \
    CGO_ENABLED=0 GOGC=off GOOS=linux \
    go build -ldflags $LDFLAGS -a -installsuffix nocgo -o snip .

FROM scratch
COPY --from=builder /go/src/github.com/derek-elliott/url-shortener/snip /
CMD ["/snip"]
