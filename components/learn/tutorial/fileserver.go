package tutorial

import (
	"fmt"
	"log"
	"net/http"
)

func SimpleFileServer(dir string, port int) {
	// create file server handler
	fs := http.FileServer(http.Dir(dir))

	// start HTTP server with `fs` as the default handler
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), fs))
}
