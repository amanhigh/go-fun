package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

//Event is a go representation of an http server-sent event
type SseEvent struct {
	Type string //SSE Type - event/data
	Data string //Actual Data
}

var (
	delim     = []byte{':', ' '}
	lineDelim = byte('\n')
)

func main() {
	/* Build a Get Request */
	if request, err := http.NewRequest("GET", "http://localhost:8080/stream", nil); err == nil {
		/* Execute Request and get handle to Event Channel */
		if eventChannel, err := fireSSERequest(request, context.Background()); err == nil {
			/* Listen to event channel for SSE Events */
			for event := range eventChannel {
				fmt.Printf("Event:%+v\n", event)
			}
		} else {
			fmt.Printf("Error Executing SSE Request:%+v\n", err)
		}
	} else {
		fmt.Printf("Error Building Request:%+v\n", err)
	}
}

/*
	Fires Given Request treating it as SSE Request & Provides a channel to listen for SSE Events.
	Context can be used to cancel listening to events before server closes stream.
*/
func fireSSERequest(request *http.Request, ctx context.Context) (eventChannel chan SseEvent, err error) {
	/* Add Header to accept streaming events */
	request.Header.Set("Accept", "text/event-stream")

	/* Make Channel to Report Events */
	eventChannel = make(chan SseEvent)
	var response *http.Response

	/* Fire Request */
	if response, err = http.DefaultClient.Do(request); err == nil {
		/* Open a Reader on Response Body */
		go liveRequestLoop(response, eventChannel, ctx)
	} else {
		fmt.Printf("Http Request Failed:%+v\n", err)
	}

	return
}

/*
	Given a response reads it and provides updates SSE Event updates on  channel provided to it.
	Context can be used to cancel listening to events before server closes stream.
*/
func liveRequestLoop(response *http.Response, eventChannel chan SseEvent, ctx context.Context) {
	defer response.Body.Close()
	br := bufio.NewReader(response.Body)
	for {
		select {
		case <-ctx.Done():
			close(eventChannel)
			fmt.Println("Context Signal Recieved Exiting")
			return
		default:
			/* Read Lines Upto Delimiter */
			if readBytes, err := br.ReadBytes(lineDelim); err == nil || err == io.EOF {

				/* Skip Lines without Content */
				if len(readBytes) < 2 {
					continue
				}
				eventChannel <- buildEvent(readBytes)

				/* Exit once Stream Closes */
				if err == io.EOF {
					fmt.Println("Stream Reading Finished")
					close(eventChannel)
					break
				}

			} else {
				fmt.Printf("Error Reading Line:%+v\n", err)
			}
		}
	}
}

/* Builds a SSE Event from read line */
func buildEvent(readBytes []byte) SseEvent {
	/* Split Actual Data & Marker Delimiter */
	splitLine := bytes.Split(readBytes, delim)
	/* Extract Data & Type */
	dataType := string(bytes.TrimSpace(splitLine[0]))
	data := string(bytes.TrimSpace(splitLine[1]))

	/* Construct Event */
	return SseEvent{Type: dataType, Data: data}
}
