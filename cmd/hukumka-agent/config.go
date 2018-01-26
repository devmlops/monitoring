package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/wwwthomson/monitoring/pkg/agent"
)

type Config struct {
	Debug  bool
	Server agent.Server
}

func OpenConfig(path string) Config {
	var config Config
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read %s", err)
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Failed to parse %s", err)
	}
	return config
}
