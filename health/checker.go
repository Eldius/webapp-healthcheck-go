package health

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

type ServiceType int

const (
	ServiceTypeTCP ServiceType = iota
	//ServiceTypeHTTP ServiceType = iota
)

type ServiceConfig interface {
	Name() string
	Type() ServiceType
	Endpoint() string
	Timeout() time.Duration
	Test() Status
}

type TCPServiceConfig struct {
	name     string
	endpoint string
	timeout  time.Duration
}

func (cfg *TCPServiceConfig) Name() string {
	return cfg.name
}
func (cfg *TCPServiceConfig) Type() ServiceType {
	return ServiceTypeTCP
}
func (cfg *TCPServiceConfig) Endpoint() string {
	return cfg.endpoint
}
func (cfg *TCPServiceConfig) Timeout() time.Duration {
	return cfg.timeout
}

func (cfg *TCPServiceConfig) Test() Status {
	if url, err := parseHost(cfg.Endpoint()); err != nil {
		return Status{
			Name:   "cartao-adesao",
			Status: ServiceStatusNOK,
			Details: map[string]string{
				"error": err.Error(),
			},
		}
	} else {
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

func BuildChecker(cfgList []ServiceConfig) http.HandlerFunc {
	return func(w ResponseWriter, r *Request) {
		
	}
}

func checkHealth(cfgList []ServiceConfig) HealthStatus {
	h := HealthStatus{
		BuildDate: "",
		Version:   "",
	}

	for _, c := range cfgList {
		h.Statuses = append(h.Statuses, c.Test())
	}

	return h
}
