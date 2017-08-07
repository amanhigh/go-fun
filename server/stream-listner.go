package main

import (
	"fmt"
	"net/http"
	"io"
	"bytes"
	"bufio"
)

//Event is a go representation of an http server-sent event
type Event struct {
	Type string //SSE Type - event/data
	Data string //Actual Data
}

var (
	delim     = []byte{':', ' '}
	lineDelim = byte('\n')
)

func main() {
	uri := "http://localhost:8080/stream"
	/* Build a Get Request */
	if request, err := http.NewRequest("GET", uri, nil); err == nil {
		/* Add Header to accept streaming events */
		request.Header.Set("Accept", "text/event-stream")

		/* Fire Request */
		if resp, err := http.DefaultClient.Do(request); err == nil {
			/* Open a Reader on Response Body */
			br := bufio.NewReader(resp.Body)
			defer resp.Body.Close()

			for {
				/* Read Lines Upto Delimiter */
				if readBytes, err := br.ReadBytes(lineDelim); err == nil || err == io.EOF {

					/* Skip Lines without Content */
					if len(readBytes) < 2 {
						continue
					}

					/* Split Actual Data & Marker Delimiter */
					splitLine := bytes.Split(readBytes, delim)

					/* Invalid Line if we don't have content so skip it */
					if len(splitLine) < 2 {
						continue
					}

					/* Extract Data & Type */
					dataType := string(bytes.TrimSpace(splitLine[0]))
					data := string(bytes.TrimSpace(splitLine[1]))

					/* Based on Event/Data Write Event */
					currEvent := &Event{Type: dataType, Data: data}
					fmt.Printf("Event:%+v\n", currEvent)

					/* Exit once Stream Closes */
					if err == io.EOF {
						fmt.Println("Stream Reading Finished")
						break
					}
				} else {
					fmt.Printf("Error Reading Line:%+v\n", err)
				}
			}
		} else {
			fmt.Printf("Error Building Request:%+v\n", err)
		}
	}

}
