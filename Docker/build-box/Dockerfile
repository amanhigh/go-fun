#This is a sample Image
FROM ubuntu
MAINTAINER coool.aman@gmail.com

#Install Basics
RUN apt-get update; apt-get install -y --force-yes \
ca-certificates openssh-client \
wget build-essential git unzip \
python3.4 python3-pip;

#Clean Apt
RUN apt-get clean
RUN rm -rf /var/lib/apt/lists/*

#Install Go1.1
ENV GO_PKG=go1.14.linux-amd64.tar.gz
RUN wget https://dl.google.com/go/${GO_PKG} --no-check-certificate
RUN tar -xvf ${GO_PKG}
RUN mv go /usr/local
RUN rm ${GO_PKG}

ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH

CMD sleep 10000

