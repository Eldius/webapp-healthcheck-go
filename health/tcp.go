package health

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"time"
)

/*
TCPServiceConfig is the checker config type for TCP checks
*/
type TCPServiceConfig struct {
	name     string
	endpoint string
	timeout  time.Duration
}

/*
Name returns the test/service name
*/
func (cfg *TCPServiceConfig) Name() string {
	return cfg.name
}

/*
Type returns the test/service type (ServiceTypeTCP)
*/
func (cfg *TCPServiceConfig) Type() ServiceType {
	return ServiceTypeTCP
}

/*
Endpoint returns the test/service endpoint
*/
func (cfg *TCPServiceConfig) Endpoint() string {
	return cfg.endpoint
}

/*
Timeout returns the test/service TCP test timeout
*/
func (cfg *TCPServiceConfig) Timeout() time.Duration {
	return cfg.timeout
}

/*
Test returns the test/service status
*/
func (cfg *TCPServiceConfig) Test() Status {
	url, err := parseHost(cfg.Endpoint())
	if err != nil {
		return Status{
			Name:   "cartao-adesao",
			Status: ServiceStatusNOK,
			Details: map[string]string{
				"error": err.Error(),
			},
		}
	}

	start := time.Now()
	conn, err := net.DialTimeout("tcp", url, cfg.Timeout())
	if err != nil {
		log.Println("Something wrong: ", err)
		return Status{
			Name:   cfg.Name(),
			Status: ServiceStatusNOK,
			Details: map[string]string{
				"cause": err.Error(),
			},
		}
	}
	defer func() {
		conn.Close()
		log.Println("Connection closed")
	}()
	return Status{
		Name:   cfg.Name(),
		Status: ServiceStatusOK,
		Details: map[string]string{
			"time": time.Since(start).String(),
		},
	}

}

func NewTCPChecker(name string, endpoint string, timeout time.Duration) ServiceConfig {
	return &TCPServiceConfig{
		name:     name,
		endpoint: endpoint,
		timeout:  timeout,
	}
}

func parseHost(target string) (host string, err error) {
	u, err := url.ParseRequestURI(target)
	if err != nil {
		return
	}
	var port = u.Port()
	if port == "" {
		if u.Scheme == "http" {
			port = "80"
		} else if u.Scheme == "https" {
			port = "443"
		}
	}
	host = fmt.Sprintf("%s:%s", u.Hostname(), port)

	return
}
