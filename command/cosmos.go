package commander

import "fmt"

func CosmosCurl(host string, startMin int, endMin int, metric string,pipe string) string {
	cosmosUrl := fmt.Sprintf("http://%v/api/query?start=%vm-ago&end=%vm-ago&m=%v", host, startMin, endMin, metric)
	PrintWhite(cosmosUrl)
	output := Jcurl(cosmosUrl, pipe)
	PrintWhite(output)
	return output
}
