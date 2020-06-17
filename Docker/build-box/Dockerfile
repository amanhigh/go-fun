#This is a sample Image
FROM debian:9
MAINTAINER coool.aman@gmail.com

#Install Basics
RUN apt-get update; apt-get install -y --force-yes ca-certificates wget build-essential openssh-client git unzip python python-pip
#Clean Apt
RUN apt-get clean
RUN rm -rf /var/lib/apt/lists/*

#Install Go1.1
RUN wget https://dl.google.com/go/go1.12.9.linux-amd64.tar.gz
RUN tar -xvf go1.12.9.linux-amd64.tar.gz
RUN mv go /usr/local
RUN rm go1.12.9.linux-amd64.tar.gz

ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH

CMD sleep 10000

