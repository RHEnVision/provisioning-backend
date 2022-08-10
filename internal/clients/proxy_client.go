package clients

import (
	"net/http"
	"net/url"
)

func GetProxiedClient() (*http.Client, error) {
	proxyStr := "<PROXY_URL>"
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Transport: transport,
	}
	return client, nil
}
