jq '.prod.remove | length' ./sample.json
jq '.prod | keys' ./sample.json
jq '.prod.remove | keys' ./sample.json
jq '.prod | keys' ./sample.json

echo -en "\033[1;32m Actions \033[0m \n"
jq -r '.prod | .[] |.[] |.[] | .action' ./sample.json
echo -en "\033[1;32m Names \033[0m \n"
jq -r '.prod | .[] |.[] |.[] | .users | .[] | .identifier' ./sample.json

echo -en "\033[1;32m Filtering \033[0m \n"
jq -r '.prod | .[] |.[] |.[] | select((.users|.[]|.identifier=="singh") and (.action=="link")) | .effect ' ./sample.json
