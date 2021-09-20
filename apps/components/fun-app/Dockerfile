# Build docker build -t <APP_NAME> .
# Eg. docker build -t fun-app .
FROM scratch
ADD ./fun /
ADD ./apps/components/fun-app/config.yml /etc/fun-app/config.yml
#Create Log Directory in scratch file
ADD ./apps/components/fun-app/config.yml /var/log/fun-app/access.log

CMD ["/fun"]