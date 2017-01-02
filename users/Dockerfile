FROM golang:1.7

RUN apt-get update && \
    apt-get install -y \
        ca-certificates \
        git \
        ssh \
        build-essential \
        golang && \
    apt-get autoremove && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

ENV GOROOT="/usr/local/go/"
ENV GOPATH="/go:$GOROOT"
ENV PATH="$PATH:/go/bin:$GOROOT/bin"
RUN mkdir -p /go/src/github.com/snickers54/microservices/login
ADD . /go/src/github.com/snickers54/microservices/login/

WORKDIR /go/src/github.com/snickers54/microservices/login
RUN go get && go build -o login
CMD ./login ./config.yaml
ENTRYPOINT ["./login"]
