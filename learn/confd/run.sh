#!/usr/bin/env bash

#File
confd -confdir . -onetime -backend file -file myapp.yaml

#Etcd
#etcdctl set /myapp/database/url db.example.com
#etcdctl set /myapp/database/user rob
#confd -confdir . -watch -backend etcd -node http://192.168.99.100:2379