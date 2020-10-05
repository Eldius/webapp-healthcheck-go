package health

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

//ServiceType service type definition
type ServiceType int

const (
	// ServiceTypeTCP is the TCP type checker
	ServiceTypeTCP ServiceType = iota
	//ServiceTypeHTTP ServiceType = iota
)

/*
ServiceConfig defines the checker interface
*/
type ServiceConfig interface {
	Name() string
	Type() ServiceType
	Endpoint() string
	Timeout() time.Duration
	Test() Status
}

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

/*
BuildChecker build the Checker responder
*/
func BuildChecker(cfgList []ServiceConfig, info map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := checkHealth(cfgList, info)
		if h.Status == ServiceStatusOK {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(502)
		}

		_ = json.NewEncoder(w).Encode(h)
	}
}

func checkHealth(cfgList []ServiceConfig, info map[string]string) HealthStatus {
	h := HealthStatus{
		Info: info,
	}

	for _, c := range cfgList {
		h.Services = append(h.Services, c.Test())
	}

	h.Evaluate()

	return h
}
