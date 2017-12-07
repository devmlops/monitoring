package agent

import (
    "time"
    "net/http"
)

type Server struct {
    IP string
    port string
}

var server = Server{IP: "127.0.0.1", port: "8080"}

var tr = &http.Transport{
	IdleConnTimeout:    2 * time.Second,
}

var client = &http.Client{Transport: tr}
