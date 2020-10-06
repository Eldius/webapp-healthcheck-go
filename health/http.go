package health

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type HTTPServiceConfig struct {
	name     string
	endpoint string
	timeout  time.Duration
	client *http.Client
	status int
}

/*
Name returns the test/service name
*/
func (cfg *HTTPServiceConfig) Name() string {
	return cfg.name
}

/*
Type returns the test/service type (CheckerTypeTCP)
*/
func (cfg *HTTPServiceConfig) Type() CheckerType {
	return CheckerTypeTCP
}

/*
Endpoint returns the test/service endpoint
*/
func (cfg *HTTPServiceConfig) Endpoint() string {
	return cfg.endpoint
}

/*
Timeout returns the test/service TCP test timeout
*/
func (cfg *HTTPServiceConfig) Timeout() time.Duration {
	return cfg.timeout
}

/*
Test returns the test/service status
*/
func (cfg *HTTPServiceConfig) Test() Status {
	start := time.Now()
	cfg.client.Timeout = cfg.timeout

	r, err := cfg.client.Get(cfg.Endpoint())

	if err != nil {
		return Status{
			Name:   cfg.Name(),
			Status: CheckerStatusNOK,
			Details: map[string]string{
				"time": time.Since(start).String(),
				"cause": err.Error(),
			},
		}
	}

	if r.StatusCode != cfg.status {
		return Status{
			Name:   cfg.Name(),
			Status: CheckerStatusNOK,
			Details: map[string]string{
				"time": time.Since(start).String(),
				"status": strconv.Itoa(r.StatusCode),
			},
		}
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return Status{
			Name:   cfg.Name(),
			Status: CheckerStatusNOK,
			Details: map[string]string{
				"time": time.Since(start).String(),
				"status": strconv.Itoa(r.StatusCode),
				"cause": err.Error(),
			},
		}
	}
	return Status{
		Name:   cfg.Name(),
		Status: CheckerStatusOK,
		Details: map[string]string{
			"time": time.Since(start).String(),
			"status": strconv.Itoa(r.StatusCode),
			"response": string(body),
		},
	}
}

func NewHTTPChecker(name string, endpoint string, timeout time.Duration, status int) ServiceChecker {
	return &HTTPServiceConfig{
		name:     name,
		endpoint: endpoint,
		timeout:  timeout,
		client: &http.Client{
			Timeout: timeout,
		},
		status: status,
	}
}