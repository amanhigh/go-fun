#This is a sample Image
# HACK: #A Migrate to alpine mini Tutorial File
FROM ubuntu
MAINTAINER name@gmail.com

RUN apt-get update
RUN apt-get install -y nginx apache2 apache2-utils curl

#Clean Apt
RUN apt-get clean
RUN rm -rf /var/lib/apt/lists/*

ENV var1=Tutorial var2=point

EXPOSE 80
WORKDIR /newTemp

#CMD echo "Container Started: $var1 $var2";apache2ctl -D FOREGROUND
CMD echo "Container Started: $var1 $var2";apache2ctl start;/bin/bash
