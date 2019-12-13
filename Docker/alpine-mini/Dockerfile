#This is a sample Image
FROM alpine
MAINTAINER coool.aman@gmail.com

RUN apk update

#add bash shell
RUN apk add --update bash
#add root certs
RUN apk add ca-certificates

RUN rm -rf /var/cache/apk/*

CMD echo "Resource Limits: -m 500m --cpus=0.05 --read-only --mount type=tmpfs,destination=/tmp,tmpfs-size=500m"
CMD echo "Small Linux with Bash";/bin/bash
