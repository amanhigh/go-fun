/opt/bitnami/zookeeper/bin/zkCli.sh -server localhost:2181 <<EOF
create /FirstNode “FirstData”
create -s  /SequentialNode "SequentialData-1"
create -s /SequentialNode "SequentialData-2"
create -e /EphimeralNode "EphimiralData"

ls /

get /FirstNode
set /FirstNode “FirstOverwrittenData”

get /FirstNode
get /EphimeralNode

stat /EphimeralNode

delete /FirstNode
delete /EphimeralNode
quit
EOF

# deleteall /
