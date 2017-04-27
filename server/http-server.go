package main

import (
	"net/http"
	"fmt"
	"encoding/json"
)

func main() {
	http.HandleFunc("/", handleRoot)
	http.ListenAndServe(":8080", nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{"aman": "Preet"}
	if result, e := json.Marshal(p); e == nil {
		fmt.Fprint(w, string(result))
	} else {
		fmt.Println("Error:", e)
	}
}
