#This is a sample Image
FROM debian:7
MAINTAINER coool.aman@gmail.com

#Install Basics
RUN apt-get update; apt-get install -y --force-yes ca-certificates wget build-essential openssh-client git

#Install Go1.1
RUN wget https://dl.google.com/go/go1.11.linux-amd64.tar.gz
RUN tar -xvf go1.11.linux-amd64.tar.gz
RUN mv go /usr/local

ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH

CMD sleep 10000