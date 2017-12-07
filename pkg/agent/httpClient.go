package agent

import (
    "time"
    "net/http"
)

var tr = &http.Transport{
	IdleConnTimeout:    2 * time.Second,
}

var client = &http.Client{Transport: tr}
