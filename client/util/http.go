package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
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
	if global.ClientConfig.Token != "" {
		req.Header.Add("Authorization", "Bearer "+global.ClientConfig.Token)
	} else {
		return nil, errors.New("no token there")
	}

directDo:
	return global.HttpClient.Do(req)
}

func Request2(method, baseUrl string, requestData any, params map[string]string, authRequired bool) (*http.Response, error) {
	var req *http.Request
	var err error
	if method == http.MethodGet {
		query := url.Values{}
		for k, v := range params {
			query.Add(k, v)
		}

		fullUrl := fmt.Sprintf("%s?%s", baseUrl, query.Encode())
		req, err = http.NewRequest(method, fullUrl, nil)
		if err != nil {
			zap.S().Errorf("failed to create a instance for GET method request: %v", err)
			return nil, err
		}
	}

	if method == http.MethodPost {
		var rd []byte
		rd, err = json.Marshal(requestData)
		if err != nil {
			return nil, err
		}

		body := bytes.NewReader(rd)

		req, err = http.NewRequest(method, baseUrl, body)
		if err != nil {
			zap.S().Errorf("failed to create a instance for POST method request: %v", err)
			return nil, err
		}
		req.Header.Add("Content-Type", "application/json")
	}

	if !authRequired {
		goto directDo
	}

	// 如果global.Token不为空，则自动添加到请求头
	if global.ClientConfig.Token != "" {
		req.Header.Add("Authorization", "Bearer "+global.ClientConfig.Token)
	} else {
		return nil, errors.New("no token there")
	}

directDo:
	return global.HttpClient.Do(req)
}
