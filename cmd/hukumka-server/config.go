package main

import (
	"io/ioutil"
	"log"
	"encoding/json"
)

type Config struct {
	Server      ServerConfig
	TelegramBot TelegramConfig
	Network     NetworkConfig
	Memory      MemoryConfig
	Swap        SwapConfig
	CPU         CpuConfig
	Disk        DiskConfig
}

type ServerConfig struct {
	Address string
}

type TelegramConfig struct {
	Token string
	Users []int64
}

type NetworkConfig struct {
	Max      uint64
	MaxLimit uint64
}

type MemoryConfig struct {
	Max      uint64
	MaxLimit uint64
}

type SwapConfig struct {
	Max      uint64
	MaxLimit uint64
}

type CpuConfig struct {
	Max      uint64
	MaxLimit uint64
}

type DiskConfig struct {
	Max      uint64
	MaxLimit uint64
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

func OpenMessage(path string) string {
	//var message string
	message, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read %s", err)
	}
	return string(message)
}
