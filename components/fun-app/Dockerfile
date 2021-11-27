# Build docker build -t <APP_NAME> .
# Eg. docker build -t fun-app .
FROM scratch
ADD ./fun /
#Create Log Directory in scratch file
ADD ./components/fun-app/fun-app.go /var/log/fun-app/access.log

CMD ["/fun"]