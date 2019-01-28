#!/usr/bin/env bash

#File
confd -confdir . -onetime -backend file -file myapp.yaml

#Etcd
#etcdctl set /myapp/database/url db.example.com
#etcdctl set /myapp/database/user rob
#confd -confdir . -watch -backend etcd -node http://docker:2379

#curl -X PUT -d 'db.example.com' http://docker:8500/v1/kv/myapp/database/url
#curl -X PUT -d 'rob' http://docker:8500/v1/kv/myapp/database/user
#confd -confdir . -watch -backend consul -node http://docker:8500