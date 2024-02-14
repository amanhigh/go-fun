echo -en "\033[1;32m Vegeta Attack (Login Host) \033[0m \n"
kubectl run vegeta -n fun-app --image="peterevans/vegeta" -- sh -c "sleep 10000"
echo -en "\033[1;33m HELM: echo 'GET http://fun-app:9000/person/all' | vegeta attack | vegeta report \033[0m \n"
echo -en "\033[1;33m DEVSPACE: echo 'GET http://app:8080/person/all' | vegeta attack | vegeta report \033[0m \n"

echo -en "\033[1;32m Metrics (Istio Only) \033[0m \n"
echo -en "\033[1;33m Prometheus: http://localhost:9090/graph?g0.expr=rate(fun_app_person_count%5B5m%5D)&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=5m&g1.expr=fun_app_person_create_time_bucket&g1.tab=0&g1.stacked=1&g1.show_exemplars=1&g1.range_input=1h&g2.expr=rate(fun_app_person_create_time_count%5B5m%5D)&g2.tab=0&g2.stacked=0&g2.show_exemplars=0&g2.range_input=1h \033[0m \n"
echo -en "\033[1;33m Grafana Import: /fun-app/it/dashboard.json \033[0m \n"