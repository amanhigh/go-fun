FROM python:3.8-slim
MAINTAINER coool.aman@gmail.com

#Install Webssh
RUN pip3 install webssh==1.5.2

CMD wssh --address=0.0.0.0 --port=8080 --fbidhttp=False

EXPOSE 8080