#Build CMD: docker build . -t 10.47.7.214/webshell:0.2
#Run CMD: docker run -it 10.47.7.214/webshell:0.2 -p 8080:8080
FROM python:3.8-slim
MAINTAINER coool.aman@gmail.com

#Install Webssh
RUN pip3 install webssh==1.5.2

CMD wssh --address=0.0.0.0 --port=8080 --fbidhttp=False

EXPOSE 8080