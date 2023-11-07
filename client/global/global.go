package global

import (
	"net/http"
	"time"
	"to-persist/client/config"
)

var (
	ClientConfig = &config.Config{}

	HttpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
)
