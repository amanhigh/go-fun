#RUN: docker build -t fun . && docker run -it fun
#This is a sample Image
FROM alpine
MAINTAINER name@gmail.com

# Set Timezone
ENV TZ=Asia/Kolkata

#Add Repo
RUN echo 'http://dl-cdn.alpinelinux.org/alpine/v3.4/main' >> /etc/apk/repositories

#add packages
RUN apk add --update bash curl ca-certificates

#Clean installtion files
RUN rm -rf /var/cache/apk/*

# Run as Custom User
ARG user=aman
ARG group=fun
ARG uid=7500
ARG gid=8500
RUN addgroup -g ${gid} ${group} \
    && adduser -u ${uid} --ingroup ${group} -D -H ${user}

RUN chown -R ${user}:${group} /usr/share/ /etc/ /var/log/
USER ${user}

CMD echo "Small Linux with Bash; Optional Resource Limits: -m 500m --cpus=0.05 --read-only --mount type=tmpfs,destination=/tmp,tmpfs-size=500m";/bin/bash
