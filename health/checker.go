package health

import (
	"encoding/json"
	"net/http"
	"time"
)

//CheckerType service type definition
type CheckerType int

const (
	// CheckerTypeTCP is the TCP type checker
	CheckerTypeTCP CheckerType = iota
	// CheckerTypeDB is the database type checker
	CheckerTypeDB CheckerType = iota
	//CheckerTypeHTTP CheckerType = iota
)

/*
ServiceChecker defines the checker interface
*/
type ServiceChecker interface {
	Name() string
	Type() CheckerType
	Timeout() time.Duration
	Test() Status
}

/*
BuildChecker build the Checker responder
*/
func BuildChecker(cfgList []ServiceChecker, info map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h := checkHealth(cfgList, info)
		if h.Status == CheckerStatusOK {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(502)
		}

		_ = json.NewEncoder(w).Encode(h)
	}
}

func checkHealth(cfgList []ServiceChecker, info map[string]string) HealthStatus {
	h := HealthStatus{
		Info: info,
	}

	for _, c := range cfgList {
		h.Services = append(h.Services, c.Test())
	}

	h.Evaluate()

	return h
}
