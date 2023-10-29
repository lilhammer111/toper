package util

import (
	"errors"
	"io"
	"net/http"
	"to-persist/client/global"
)

func Request(method, url string, body io.Reader, authRequired bool) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if method == http.MethodPost {
		req.Header.Add("Content-Type", "application/json")
	}

	if !authRequired {
		goto directDo
	}

	// 如果global.Token不为空，则自动添加到请求头
	if global.Config.Token != "" {
		req.Header.Add("Authorization", "Bearer "+global.Config.Token)
	} else {
		return nil, errors.New("no token there")
	}

directDo:
	return global.HttpClient.Do(req)
}
