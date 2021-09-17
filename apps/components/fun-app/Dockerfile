# Build docker build -t <APP_NAME> .
# Eg. docker build -t fun-app .
FROM scratch
ADD ./fun /
ADD ./apps/components/fun-app/config.yml /etc/fun-app/config.yml
CMD ["/fun"]