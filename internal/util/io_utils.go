package util

import (
	"net/http"
	"time"
)

var (
	DimHTTPClient = http.Client{
		Timeout: time.Second * 5,
	}
)
