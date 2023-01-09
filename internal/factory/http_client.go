package factory

import (
	"github.com/go-resty/resty/v2"
	"github.com/kozmod/progen/internal/config"
)

func NewHTTPClient(conf *config.HTTPClient) *resty.Client {
	if conf == nil {
		return resty.New()
	}

	client := resty.New().
		SetHeaders(conf.Headers).
		SetBaseURL(conf.BaseURL.String())

	client.Debug = conf.Debug
	return client
}
