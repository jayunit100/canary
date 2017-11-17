package main

import (
	"os"
	"testing"
)

func TestLookup(t *testing.T) {
	cfg := ViperLoad()
	if len(cfg.Services) < 4 {
		t.Fail()
	}
	for s, _ := range cfg.Services {
		t.Log(s)
	}
	if cfg.Services["zookeeper"] != 2181 {
		t.Fail()
	}
}

// make sure that services are parsed correctly
func TestLookupEnv(t *testing.T) {
	os.Setenv("ENV_CONFIG_JSON", "{\"services\":{\"blackduckabc\":1234}, \"svcTimeout\":10}")
	cfg := ViperLoad()
	if len(cfg.Services) != 1 {
		t.Fail()
	}
	if cfg.Services["blackduckabc"] != 1234 {
		t.Fail()
	}
	if cfg.SvcTimeout != 10 {
		t.Fail()
	}
}
