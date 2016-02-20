FROM ubuntu:14.04

ENV DEBIAN_FRONTEND noninteractive

RUN sudo apt-get update && sudo apt-get install -y \
    vim \
    git \
    mercurial \
    golang

# Setup golang enviroment

RUN mkdir /root/go
ENV GOPATH=/root/go
ENV PATH=$PATH:$GOPATH/bin

# Clone nebula repo and install it

WORKDIR $GOPATH/src/github.com/lizhongz/nebula
COPY . $GOPATH/src/github.com/lizhongz/nebula

RUN go get
