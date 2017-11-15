package sidecar

import (
	"testing"
)

func TestLookup(t *testing.T) {
	cfg := &Config{
		Services: map[string]int{
			"web-1": 8080,
			"web-2": 8080,
			"web-3": 8080,
		},
	}
	summary, _ := cfg.LookupHub()

	if summary["services_summary"] == "" {
		t.Fail()
	}

	web1 := summary["services_detail_web-1"]
	if web1 == "" {
		t.Log("FAIL!!!!!!!!! web1 is empty !")
		t.Fail()
	}
	if len(serviceMetrics) != 3 {
		t.Log("service map length is wrong!", len(serviceMetrics), "\n\tsummary map entries", len(serviceMetrics))
		for k, v := range serviceMetrics {
			t.Log("key:", k)
			t.Log("value:", v)
		}
		t.Fail()
	}
}
