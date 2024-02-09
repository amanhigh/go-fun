# https://stedolan.github.io/jq/manual/
echo -en "\033[1;32m Map functions \033[0m \n"
jq '.prod.remove | length' ./sample.json
jq '.prod | keys' ./sample.json
jq '.prod.remove | keys' ./sample.json
jq '.prod | keys' ./sample.json

echo -en "\033[1;32m Filtering \033[0m \n"
jq -r '.prod | .[] |.[] |.[] | .action' ./sample.json
jq -r '.prod | .[] |.[] |.[] | .users | .[] | .identifier' ./sample.json
jq -r '.prod | .[] |.[] |.[] | select((.users|.[]|.identifier=="singh") and (.action=="link")) | .effect ' ./sample.json

echo -en "\033[1;32m Replacement \033[0m \n"
cat payload.json | jq '.|.age=22'

echo -en "\033[1;32m Loops and String Interploation \033[0m \n"
echo 1 | jq -cM 'while(.<5; .+1) | {method: "POST", url: "http://:6060", body: {id: .,name:"Aman-\(.)","age": (44+.), gender: "MALE"} }'
jq -ncM --slurpfile a payload.json 'while(.<5; .+1) as $i | $a | $a[0].name="\($a[0].name)-\($i)"'

echo -en "\033[1;32m File Reading & Variables \033[0m \n"
jq -n --slurpfile a payload.json '$a[0]'
jq -ncM --slurpfile a sample.json '$a[]|keys'
jq '.prod.remove | length as $len | $len+1' ./sample.json
jq -n  'env.PWD'