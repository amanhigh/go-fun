package commander

import (
	"fmt"
	"encoding/json"
	"strconv"
	"sort"
)

func CosmosCurl(host string, startMin int, endMin int, metric string,pipe string) string {
	cosmosUrl := fmt.Sprintf("http://%v/api/query?start=%vm-ago&end=%vm-ago&m=%v", host, startMin, endMin, metric)
	PrintWhite(cosmosUrl)
	output := Jcurl(cosmosUrl, pipe)
	PrintWhite(output)
	return output
}

/*
Computes Rates from Cosmos Query which returs a running counter values
as its Dps.
*/
func CosmosRates(host string, startMin int, endMin int, metric string) []int {
	cosmosUrl := fmt.Sprintf("http://%v/api/query?start=%vm-ago&end=%vm-ago&m=%v", host, startMin, endMin, metric)
	PrintWhite(cosmosUrl)
	result := Jcurl(cosmosUrl, "jq -r '.[] | .dps'")

	/* Unmarshal Json */
	rate := map[string]int{}
	json.Unmarshal([]byte(result), &rate)

	/* Convert Timestamps to Int */
	timeStampMap := map[int]int{}
	sortedTimeStamp := []int{}
	for timeStamp, dps := range rate {
		intTimestamp, _ := strconv.Atoi(timeStamp)
		sortedTimeStamp = append(sortedTimeStamp, intTimestamp)
		timeStampMap[intTimestamp] = dps
	}

	/* Sort Time Stamps & Dps Accordingly  */
	sort.Ints(sortedTimeStamp)
	sortedDps := []int{}
	//lastTimeStamp := 0
	for _, timestamp := range sortedTimeStamp {
		//fmt.Printf("Timestamp:%v Last:%v Diff:%v\n", timestamp, lastTimeStamp, timestamp-lastTimeStamp)
		sortedDps = append(sortedDps, timeStampMap[timestamp])
		//lastTimeStamp = timestamp
	}

	/* Compute Rates */
	lastDps := 0
	computedRates := []int{}
	for i, dps := range sortedDps {
		//commander.PrintWhite(fmt.Sprintf("I: %v Dps: %v Last:%v Diff:%v", i, dps, lastDps, dps-lastDps))
		if i > 0 {
			computedRate := dps - lastDps
			computedRates = append(computedRates, computedRate)
		}
		lastDps = dps
	}
	return computedRates
}

