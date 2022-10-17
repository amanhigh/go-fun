etcdctl version
# etcdctl auth disable --user="root:${ETCD_ROOT_PASSWORD}"
etcdctl --user root:$ETCD_ROOT_PASSWORD put /message Hello
etcdctl --user root:$ETCD_ROOT_PASSWORD get /message