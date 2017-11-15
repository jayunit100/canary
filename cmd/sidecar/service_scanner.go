package sidecar

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Initialize the metric to 0, since at start, we have successfully seen each service 0 times.
func init() {
	viperLoad()
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		viperLoad()
	})
}

type Config struct {
	Services   map[string]int
	SvcTimeout int
}

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	serviceCheck = prometheus.NewGauge(prometheus.CounterOpts{
		Namespace: "sidecar",
		Subsystem: "metrics",
		Name:      "dns_lookup",
		Help:      "The current lookup time for a service.",
	})
)

func main() {
	http.Handle("/metrics", prometheus.Handler())

	indexed.Inc()
	size.Set(5)

	http.ListenAndServe(":8080", nil)
}

func init() {
	prometheus.MustRegister(indexed)
	prometheus.MustRegister(size)
}


func (cfg *Config) Curl(svc string, port int) (int64, error) {
	SvcTimeout := time.Duration(cfg.SvcTimeout) * time.Second
	client := http.Client{
		Timeout: SvcTimeout,
	}
	start := time.Now()
	url := strings.Join([]string{"http://" + svc, fmt.Sprintf("%v", port)}, ":")
	resp, err := client.Get(url)
	log.Info(fmt.Sprintf("Response from %v was %v : %v", url, resp, err))
	elapsed := int64(time.Since(start) / time.Millisecond)
	return elapsed, err
}



// Looks up all the services and returns a map summary.
func (cfg *Config) LookupHub() (map[string]string, error) {
	// dnslookup services status information
	srvStatus := make(map[string]string)
	// Count total failures for summary data.
	srvFailures := 0

	for service, port := range cfg.Services {

		start := time.Now()
		ip, err := net.LookupHost(service)
		elapsed := int64(time.Since(start) / time.Millisecond)

		if err != nil {
			serviceMetricsTimeMS[service] = -1
		} else {
			serviceMetricsTimeMS[service] = elapsed
		}
		if err != nil {
			srvFailures += 1
		} else {
			serviceMetrics[service] += 1
		}
		statusForThisService := fmt.Sprintf("%v (ip %v)", ip, err)
		srvStatus[service] = statusForThisService

		curlTimeMs, err := cfg.Curl(service, port)
		if err != nil {
			serviceMetrics[fmt.Sprintf("%v_curl_error", service)] += 1
		} else {
			serviceMetricsTimeMS[fmt.Sprintf("%v_curl)", service)] = curlTimeMs
		}
	}

	// Top level data structure that we return w/ all sorts of metadata in it.
	status := make(map[string]string)
	if srvFailures == len(cfg.Services) {
		status["services_summary"] = "No cfg.Services are resolving... Your cluster networking is possibly failing."
	} else if srvFailures > 0 {
		status["services_summary"] = fmt.Sprintf("Some cfg.Services are resolving ... Total failures %v.", srvFailures)
	} else {
		status["services_summary"] = "All hub cfg.Services are resolvable."
	}
	for k, v := range srvStatus {
		status["services_detail_"+k] = v
	}
	return status, nil
}

func (c *Config) url(s string) string {
	return strings.Join(
		[]string{s, string(c.Services[s])},
		":")
}

func viperLoad() *Config {
	// Default config: The blackducksoftware:hub services.  export ENV_CONFIG_JSON to override this.
	sidecar_targets := `{
	  "services":{
			"zookeeper":2181,
		  "cfssl":0,
		  "postgres":0,
		  "webapp":0,
		  "solr":0,
		  "documentation": 0
		},
		"svcTimeout":10
	}`

	if v, ok := os.LookupEnv("ENV_CONFIG_JSON"); ok {
		sidecar_targets = v
	} else {
		log.Warn(`
      ENV_CONFIG_JSON services not provided as env var
		  Instead, writing default config to sidecar.json.
		  Edit it to reload the sidecar or restart w/ the right env var.
			EXAMPLE:
				export ENV_CONFIG_JSON="{\"services\":{\"zookeeper\":2181,\"cfssl\":5555,\"postgres\":5432, \"webapp\":8080, \"solr\":0, \"documentation\": 0 }, \"svcTimeout\":10}"
      `)
	}

	d1 := []byte(sidecar_targets)
	err := ioutil.WriteFile("../../sidecar.json", d1, 0777)

	if err != nil {
		panic(fmt.Sprintf("Error writing default config file !", err))
	}
	viper.SetConfigName("sidecar") // name of config file (without extension)
	viper.AddConfigPath("./")      // path to look for the config file in

	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Errorf("Fatal error config file: %v \n", err)
	}
	var cfg *Config
	err = viper.Unmarshal(&cfg)
	return cfg
}

var hasRun bool

var serviceMetrics2 = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "request_duration_milliseconds",
		Help:    "Request latency distribution",
		Buckets: prometheus.ExponentialBuckets(10.0, 1.13, 40),
	})

var serviceMetrics map[string]int = make(map[string]int)
var serviceMetricsTimeMS map[string]int64 = make(map[string]int64)

var panics int = 0
var tries int = 0
