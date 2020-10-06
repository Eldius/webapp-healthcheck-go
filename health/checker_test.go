package health

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	defaultDuration = time.Duration(1 * time.Second)
)

func TestCheckerEndpoint(t *testing.T) {
	service0 := startservice(7777, t)
	defer service0.Close()

	service1 := startservice(8888, t)
	defer service1.Close()

	h := BuildChecker([]ServiceChecker{
		NewTCPChecker(
			"test-server0",
			"http://localhost:7777",
			defaultDuration,
		),
		NewTCPChecker(
			"test-server1",
			"http://localhost:8888",
			defaultDuration,
		),
	},
		map[string]string{
			"version":   "0.1.2",
			"buildDate": "2020-10-04 23:21:00Z",
		},
	)

	s := httptest.NewServer(h)
	defer s.Close()

	c := http.DefaultClient

	r, err := c.Get(s.URL)
	if err != nil {
		t.Errorf("Failed to make get request to test server\n%s\n", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		t.Errorf("Status code must be 200, but was %d", r.StatusCode)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Failed to read response body.\n%s\n", err.Error())
	}
	t.Logf("---\nresponse code: %s\nresponse body:\n%s\n---", r.Status, string(body))
	var hs HealthStatus
	err = json.Unmarshal(body, &hs)
	if err != nil {
		t.Errorf("Failed to unmarshal healthcheck response\n%s", err.Error())
	}

	if hs.Status != CheckerStatusOK {
		t.Errorf("Status should be UP, but was %s", hs.Status)
	}
}

func TestCheckerEndpointFail(t *testing.T) {
	service0 := startservice(7777, t)
	defer service0.Close()

	h := BuildChecker([]ServiceChecker{
		&TCPServiceConfig{
			endpoint: "http://localhost:7777",
			name:     "test-server0",
			timeout:  defaultDuration,
		},
		&TCPServiceConfig{
			endpoint: "http://localhost:8888",
			name:     "test-server1",
			timeout:  defaultDuration,
		},
	},
		map[string]string{
			"version":   "0.1.2",
			"buildDate": "2020-10-04 23:21:00Z",
		},
	)

	s := httptest.NewServer(h)
	defer s.Close()

	c := http.DefaultClient

	r, err := c.Get(s.URL)
	if err != nil {
		t.Errorf("Failed to make get request to test server\n%s\n", err.Error())
	}
	defer r.Body.Close()

	if r.StatusCode != 502 {
		t.Errorf("Status code must be 502, but was %d", r.StatusCode)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("Failed to read response body.\n%s\n", err.Error())
	}
	t.Logf("---\nresponse code: %s\nresponse body:\n%s\n---", r.Status, string(body))
	var hs HealthStatus
	err = json.Unmarshal(body, &hs)
	if err != nil {
		t.Errorf("Failed to unmarshal healthcheck response\n%s", err.Error())
	}

	if hs.Status != CheckerStatusNOK {
		t.Errorf("Status should be DOWN, but was %s", hs.Status)
	}
}

func startservice(port int, t *testing.T) net.Listener {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatal(err)
	}
	return l
}
