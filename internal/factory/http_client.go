package factory

import (
	"github.com/go-resty/resty/v2"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
)

func NewHTTPClient(conf *config.HTTPClient, logger entity.Logger) *resty.Client {
	if conf == nil {
		return resty.New()
	}

	client := resty.New().
		SetHeaders(conf.Headers).
		SetQueryParams(conf.QueryParams).
		SetBaseURL(conf.BaseURL.String())
	if logger != nil {
		client.SetLogger(logger)
	}

	client.Debug = conf.Debug
	return client
}
