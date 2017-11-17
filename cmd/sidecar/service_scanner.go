package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var registered = false

// Initialize the metric to 0, since at start, we have successfully seen each service 0 times.
func init() {
	prometheus.Unregister(prometheus.NewProcessCollector(os.Getpid(), ""))
	prometheus.Unregister(prometheus.NewGoCollector())

	cfg = ViperLoad()
	viper.WatchConfig()

	svcNames := func() []string {
		return []string{"service", "port"}
	}
	serviceCheck = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "sidecar",
			Subsystem: "metrics",
			Name:      "dns_lookup",
			Help:      "The current lookup time for a service.",
			Buckets:   []float64{1, 5, 10},
		},
		svcNames())

	prometheus.MustRegister(serviceCheck)

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		cfg = ViperLoad()
	})

	go func() {
		for {
			for svc, p := range cfg.Services {
				milliseconds, err := cfg.Curl(svc, p)
				if err != nil {
					serviceCheck.WithLabelValues(svc, fmt.Sprintf("%v", p)).Observe(float64(9999999))
				} else {
					serviceCheck.WithLabelValues(svc, fmt.Sprintf("%v", p)).Observe(float64(milliseconds))
				}
			}
			time.Sleep(5 * time.Second)
		}
	}()

}

type Config struct {
	Services   map[string]int
	SvcTimeout int
}

var (
	serviceCheck *prometheus.HistogramVec
	cfg          *Config
)

func main() {
	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Shutdown now!\n")
		os.Exit(0)
	})
	log.Info("Serving")
	http.ListenAndServe(":3000", nil)
	log.Info("Server started !")

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

func (c *Config) url(s string) string {
	return strings.Join(
		[]string{s, string(c.Services[s])},
		":")
}

var hasRun bool

func ViperLoad() *Config {
	// Default config: The blackducksoftware:hub services.  export ENV_CONFIG_JSON to override this.
	sidecarTargets := `{
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
		sidecarTargets = v
	} else {
		log.Warn(`
      ENV_CONFIG_JSON services not provided as env var
		  Instead, writing default config to sidecar.json.
		  Edit it to reload the sidecar or restart w/ the right env var.
			EXAMPLE:
				export ENV_CONFIG_JSON="{\"services\":{\"zookeeper\":2181,\"cfssl\":5555,\"postgres\":5432, \"webapp\":8080, \"solr\":0, \"documentation\": 0 }, \"svcTimeout\":10}"
      `)
	}
	d1 := []byte(sidecarTargets)
	err := ioutil.WriteFile("../../sidecar.json", d1, 0777)
	if err != nil {
		panic(fmt.Sprintf("Error writing default config file !", err))
	}

	viper.SetConfigName("sidecar") // name of config file (without extension)
	viper.AddConfigPath("../../")  // path to look for the config file in
	err = viper.ReadInConfig()     // Find and read the config file
	if err != nil {
		log.Errorf("Fatal error config file: %v \n", err)
	}

	var cfg *Config
	err = viper.Unmarshal(&cfg)

	return cfg
}
