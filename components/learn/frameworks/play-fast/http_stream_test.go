package play_fast

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

/* Server */
func handleRoot(w http.ResponseWriter, r *http.Request) {
	p := map[string]string{"aman": "Preet"}
	if result, e := json.Marshal(p); e == nil {
		fmt.Fprint(w, string(result))
	} else {
		log.Error().Err(e).Msg("Root Handler")
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
			fmt.Fprintf(w, "data: Streaming - %v\n\n", i)
			f.Flush()
			time.Sleep(time.Second)
		}
	}

	fmt.Fprint(w, "data: Streaming Finished :)\n\n")
}

/* Client */
// Event is a go representation of an http server-sent event
type SseEvent struct {
	Type string //SSE Type - event/data
	Data string //Actual Data
}

var (
	delim     = []byte{':', ' '}
	lineDelim = byte('\n')
)

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
		log.Error().Err(err).Msg("Http SSE Request")
	}

	return
}

/*
Given a response reads it and provides updates SSE Event updates on  channel provided to it.
Context can be used to cancel listening to events before server closes stream.
*/
func liveRequestLoop(response *http.Response, eventChannel chan SseEvent, ctx context.Context) {
	defer response.Body.Close()
	defer close(eventChannel)

	br := bufio.NewReader(response.Body)
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("Context Signal Recieved Exiting")
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
					log.Debug().Msg("Stream Closed")
					return
				}
			} else {
				log.Error().Err(err).Msg("Read Line Error")
				return
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

var _ = Describe("HttpStream", func() {
	var (
		srv          *http.Server
		port         = 7080
		url          = fmt.Sprintf("http://localhost:%d/stream", port)
		err          error
		request      *http.Request
		eventChannel chan SseEvent
	)

	BeforeEach(func() {
		http.HandleFunc("/", handleRoot)
		http.HandleFunc("/stream", stream)
		srv = util.NewTestServer(fmt.Sprintf(":%d", port))
		go srv.ListenAndServe() //nolint:errcheck
	})

	AfterEach(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	})

	It("should stream events", func() {
		/* Build a Get Request */
		request, err = http.NewRequest("GET", url, nil)
		Expect(err).To(BeNil())

		/* Execute Request and get handle to Event Channel */
		eventChannel, err = fireSSERequest(request, context.Background())
		Expect(err).To(BeNil())

		/* Listen to event channel for SSE Events */
		Eventually(eventChannel, 100*time.Millisecond).Should(Receive())
	})

})
