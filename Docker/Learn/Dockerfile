#This is a sample Image
FROM ubuntu
MAINTAINER coool.aman@gmail.com

RUN apt-get update
RUN apt-get install -y nginx apache2 apache2-utils curl
RUN apt-get clean

ENV var1=Tutorial var2=point

EXPOSE 80
WORKDIR /newTemp

#CMD echo "Container Started: $var1 $var2";apache2ctl -D FOREGROUND
CMD echo "Container Started: $var1 $var2";apache2ctl start;/bin/bash
