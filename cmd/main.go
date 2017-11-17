package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"github.com/blackducksoftware/canary"
)

var Version = "<None: Please build w/ 'go build ldflags -X sidecar.buildstamp `date` -X sidecar.gitinfo `git status`'>"

func main() {

	// Note that mux will parse { var } variables for free for us.
	r := mux.NewRouter()

	log.Info("Router created ")

	// Allow / or /info for simplicity to get the major data.
	r.Path("/status").Methods("GET").HandlerFunc(sidecar.StatusHandler)
	r.Path("/env").Methods("GET").HandlerFunc(sidecar.EnvHandler)


	log.Info("===========  Now launching negroni...this might take a second...")
	n := negroni.Classic()
	n.UseHandler(r)

  // nice if we could format the damn json
  // f := negroni.TextPanicFormatter
	n.Run(":3000")

	log.Info("Done ! Web app is now running.")
}
