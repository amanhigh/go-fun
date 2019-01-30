#!/usr/bin/env bash

#File
confd -confdir . --config-file confd.toml

#Etcd
#etcdctl set /myapp/database/url db.example.com
#etcdctl set /myapp/database/user rob
#confd -confdir . --config-file confd-etcd.toml

#curl -X PUT -d 'db.example.com' http://docker:8500/v1/kv/myapp/database/url
#curl -X PUT -d 'rob' http://docker:8500/v1/kv/myapp/database/user
#confd -confdir . --config-file confd-consul.toml