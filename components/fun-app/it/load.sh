# kubectl run vegeta --image="peterevans/vegeta" -- sh -c "sleep 10000"
# Attach and Run: echo "GET http://app:8080/person/all" | vegeta attack | vegeta report (dev- app:8080, image- fun-app:9000)

OP=`echo "read write all" | tr ' ' '\n' | gum filter --prompt "Load Type: "`
DURATION=`gum input --prompt "Duration: " --value=30s`
export URL=http://localhost:9000

echo -en "\033[1;32m Running $OP for $DURATION \033[0m \n";
case $OP in
  read)
     jq -ncM 'while(.<1000; .+1) | {method: "GET", url: "\(env.URL)/person/Aman-\(.)"}' | vegeta attack -format=json -duration=$DURATION | vegeta report
    ;;
  
  all)
    echo "GET $URL/person/all" | vegeta attack -duration=$DURATION | vegeta report
    ;;

  write)
    jq -ncM --argfile a payload.json \
    'while(true; .+1) as $i | $a | $a.name="\($a.name)-\($i)" | {method: "POST", header: {"Content-Type": ["application/json"]}, url: "\(env.URL)/person", body: .| @base64}' | \
    vegeta attack -lazy -format=json -duration=$DURATION | vegeta report
    ;;

  *)
    echo "Invalid OP"
    ;;
esac

