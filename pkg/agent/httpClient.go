package agent

import (
	"net/http"
	"time"
)

type Server struct {
	IP   string
	Port string
}

var tr = &http.Transport{
	IdleConnTimeout: 2 * time.Second,
}

var client = &http.Client{Transport: tr}
