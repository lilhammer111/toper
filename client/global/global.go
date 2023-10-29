package global

import (
	"net/http"
	"time"
	"to-persist/client/config"
)

var (
	Config     = &config.Config{}
	HttpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
)
