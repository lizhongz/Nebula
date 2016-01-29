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

# Clone nebula repo and install it

RUN go get -d github.com/lizhongz/nebula

ENV NE_PATH=$GOPATH/src/github.com/lizhongz/nebula

RUN cd $NE_PATH && \
    git fetch && \
    git checkout dev
