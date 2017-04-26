package main

import (
	"fmt"
	"encoding/json"
)

type person struct {
	Name         string `json:"name"`
	Age          int
	MobileNumber int64
}

func main() {
	p1 := person{"Aman", 29, 9844415553}
	//fmt.Println(p1)
	if result, e := json.Marshal(p1); e == nil {
		fmt.Println("Encoded:", string(result))
		var pDecoded person
		json.Unmarshal(result, &pDecoded)
		fmt.Printf("Decoded: %+v\n", pDecoded)
	} else {
		fmt.Println("Error:", e)
	}
}
