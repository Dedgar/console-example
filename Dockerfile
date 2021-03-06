FROM golang:latest 

ENV GOPATH=/go

ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

COPY scripts/ /usr/local/bin/

RUN mkdir -p /go/src/github.com/dedgar/console-example \
             /cert
WORKDIR /go/src/github.com/dedgar/console-example

COPY . /go/src/github.com/dedgar/console-example
RUN go-wrapper download && \
    go-wrapper install

EXPOSE 8443

USER 1001

CMD ["/usr/local/bin/start.sh"]
