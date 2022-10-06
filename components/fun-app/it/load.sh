OP=`echo "read write" | tr ' ' '\n' | gum filter --prompt "Load Type: "`
DURATION=`gum input --prompt "Duration: " --value=30s`
export URL=http://localhost:9000/person

echo -en "\033[1;32m Running $OP for $DURATION \033[0m \n";
case $OP in
  read)
    echo "GET http://localhost:9000/person/all" | vegeta attack -duration=$DURATION | vegeta report
    ;;

  write)
    jq -ncM --argfile a payload.json \
    'while(true; .+1) as $i | $a | $a.name="\($a.name)-\($i)" | {method: "POST", header: {"Content-Type": ["application/json"]}, url: env.URL, body: .| @base64}' | \
    vegeta attack -lazy -format=json -duration=$DURATION | vegeta report
    ;;

  *)
    echo "Invalid OP"
    ;;
esac

