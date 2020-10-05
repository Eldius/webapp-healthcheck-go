package health

import (
	"encoding/json"
	"net/http"
	"time"
)

//ServiceType service type definition
type ServiceType int

const (
	// ServiceTypeTCP is the TCP type checker
	ServiceTypeTCP ServiceType = iota
	// ServiceTypeDB is the database type checker
	ServiceTypeDB ServiceType = iota
	//ServiceTypeHTTP ServiceType = iota
)

/*
ServiceConfig defines the checker interface
*/
type ServiceConfig interface {
	Name() string
	Type() ServiceType
	Timeout() time.Duration
	Test() Status
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
