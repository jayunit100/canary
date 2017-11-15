/*
* Welcome to hub sidecar!
* This swims with the other hub container ducks and keeps an eye on them.
 */
package sidecar

import (
	"os"
	//	"bytes"
	"net/http"

	"strings"

	"encoding/json"

	"fmt"

	log "github.com/sirupsen/logrus"
)

func StatusHandler(rw http.ResponseWriter, req *http.Request) {
	cfg := viperLoad()

	log.Info("StatusHandler")
	promMap, err := cfg.LookupHub()

	tries = tries + 1
	fmt.Fprintf(rw, "nslookup_tries %v \n", tries)

	if err != nil {
		panics = panics + 1
	}
	fmt.Fprintf(rw, "\n# warning! if non-zero, there was a low level failure in NS lookup !\n")
	fmt.Fprintf(rw, "panics %v \n", panics)

	for k, v := range promMap {
		fmt.Fprintf(rw, "\n# %v = %v \n", k, v)
	}

	fmt.Fprintf(rw, "\n# service metrics: total lookups\n")
	for s, c := range serviceMetrics {
		fmt.Fprintf(rw, "%v_total  %v \n", s, c)
	}

	fmt.Fprintf(rw, "\n# service lookup millisecond timings (-1 == failed ) \n")

	for s, c := range serviceMetricsTimeMS {
		fmt.Fprintf(rw, "%v_time %v \n", s, c)
	}
}

func EnvHandler(rw http.ResponseWriter, req *http.Request) {
	log.Info("EnvHandler")

	environment := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		key := splits[0]
		val := strings.Join(splits[1:], "=")
		environment[key] = val
	}

	envJSON := HandleError(json.MarshalIndent(environment, "", "  ")).([]byte)
	rw.Write(envJSON)
}

func HandleError(result interface{}, err error) (r interface{}) {
	if err != nil {
		print("ERROR :  " + err.Error())
	}
	return result
}
