package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"time"
)

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/stream", stream)
	fmt.Printf("Listening on :%+v\n", "localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error:%+v\n", err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{"aman": "Preet"}
	if result, e := json.Marshal(p); e == nil {
		fmt.Fprint(w, string(result))
	} else {
		fmt.Println("Error:", e)
	}
}

func stream(w http.ResponseWriter, r *http.Request) {
	// Set the headers related to event streaming.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Listen to the closing of the http connection via the CloseNotifier
	//notify := w.(http.CloseNotifier).CloseNotify()

	if f, ok := w.(http.Flusher); !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	} else {
		for i := 0; i < 10; i++ {
			fmt.Fprint(w, fmt.Sprintf("Streaming: %v\n", i))
			f.Flush()
			time.Sleep(time.Second)
		}
	}

	fmt.Fprint(w, "Streaming Finished :)")

}
