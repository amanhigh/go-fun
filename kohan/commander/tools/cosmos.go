package tools

import (
	"fmt"
	"encoding/json"
	"strconv"
	"sort"
	"github.com/amanhigh/go-fun/kohan/commander"
	"github.com/amanhigh/go-fun/util"
)

func CosmosCurl(host string, startMin int, endMin int, metric string, pipe string) string {
	cosmosUrl := fmt.Sprintf("http://%v/api/query?start=%vm-ago&end=%vm-ago&m=%v", host, startMin, endMin, metric)
	if commander.IsDebugMode() {
		util.PrintWhite(cosmosUrl)
	}
	output := Jcurl(cosmosUrl, pipe)
	return output
}

/*
Computes Rates from Cosmos Query which returs a running counter values
as its Dps.
*/
func CosmosRates(host string, startMin int, endMin int, metric string) []int {
	cosmosUrl := fmt.Sprintf("http://%v/api/query?start=%vm-ago&end=%vm-ago&m=%v", host, startMin, endMin, metric)
	if commander.IsDebugMode() {
		util.PrintWhite(cosmosUrl)
	}
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

	/* Sort Time Stamps & Compute Rates  */
	sort.Ints(sortedTimeStamp)
	computedRates := []int{}
	lastTimeStamp := 0
	lastDps := 0
	for i, timestamp := range sortedTimeStamp {
		dps := timeStampMap[timestamp]
		timeStampDiff := timestamp - lastTimeStamp
		dpsDiff := dps - lastDps
		computedRate := dpsDiff / timeStampDiff
		if i > 0 {
			//fmt.Printf("I: %v Timestamp:%v LastTime:%v TimeDiff:%v Dps: %v LastDps:%v DpsDiff:%v Rate:%v\n", i, timestamp, lastTimeStamp, timeStampDiff, dps, lastDps, dpsDiff, computedRate)
			computedRates = append(computedRates, computedRate)
		}
		lastTimeStamp = timestamp
		lastDps = dps
	}

	return computedRates
}
