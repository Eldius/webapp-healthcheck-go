package health

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	defaultDuration = time.Duration(1 * time.Second)
)

func TestParseHostHttp(t *testing.T) {
	h, err := parseHost("http://any-host/abc")
	if err != nil {
		t.Errorf("Nao deve retornar erro, mas retornou o seguinte\n%s", err.Error())
	}

	if h != "any-host:80" {
		t.Errorf("Deveria retornar o host 'any-host:80', mas retornou '%s'", h)
	}
}

func TestParseHostHttps(t *testing.T) {
	h, err := parseHost("https://any-host/abc")
	if err != nil {
		t.Errorf("Nao deve retornar erro, mas retornou o seguinte\n%s", err.Error())
	}

	if h != "any-host:443" {
		t.Errorf("Deveria retornar o host 'any-host:443', mas retornou '%s'", h)
	}
}

func TestParseHostHttpCustomPort(t *testing.T) {
	h, err := parseHost("https://any-host:9999/abc")
	if err != nil {
		t.Errorf("Nao deve retornar erro, mas retornou o seguinte\n%s", err.Error())
	}

	if h != "any-host:9999" {
		t.Errorf("Deveria retornar o host 'any-host:9999', mas retornou '%s'", h)
	}
}

func TestParseHostTargetInvalido(t *testing.T) {
	h, err := parseHost("any-host")
	if err == nil {
		t.Error("Deve retornar erro, mas nao retornou o seguinte")
	}

	if h != "" {
		t.Errorf("Deveria retornar o host '', mas retornou '%s'", h)
	}
}

func TestTcpTestOK(t *testing.T) {
	service := startservice(9999, t)
	defer service.Close()
	cfg := TCPServiceConfig{
		name:     "success-test",
		endpoint: "http://localhost:9999",
		timeout:  defaultDuration,
	}

	s := cfg.Test()

	log.Println(s.Details)
	if s.Status != ServiceStatusOK {
		t.Errorf("Status should be UP, but was '%s'", s.Status)
	}
}

func TestTcpTestHostUnavailable(t *testing.T) {
	cfg := TCPServiceConfig{
		name:     "invalid-test",
		endpoint: "http://abc.xyz:1234",
		timeout:  defaultDuration,
	}

	s := cfg.Test()
	if s.Status != ServiceStatusNOK {
		t.Errorf("Status deve estar DOWN, mas esta '%s'", s.Status)
	}
}

func TestTcpTestTargetInvalido(t *testing.T) {
	cfg := TCPServiceConfig{
		name:     "abc.com_invalid-test",
		endpoint: "abc.com",
		timeout:  defaultDuration,
	}
	s := cfg.Test()

	log.Println(s.Details)
	if s.Status != ServiceStatusNOK {
		t.Errorf("Status deve estar DOWN, mas esta '%s'", s.Status)
	}
}

func TestCheckerEndpoint(t *testing.T)  {
	service0 := startservice(7777, t)
	defer service0.Close()

	service1 := startservice(8888, t)
	defer service1.Close()

	h := BuildChecker([]ServiceConfig{
			&TCPServiceConfig{
				endpoint: "http://localhost:7777",
				name: "test-server0",
				timeout: defaultDuration,
			},
			&TCPServiceConfig{
				endpoint: "http://localhost:8888",
				name: "test-server1",
				timeout: defaultDuration,
			},
		},
		map[string]string{
			"version": "0.1.2",
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

	if hs.Status != ServiceStatusOK {
		t.Errorf("Status should be UP, but was %s", hs.Status)
	}
}




func TestCheckerEndpointFail(t *testing.T)  {
	service0 := startservice(7777, t)
	defer service0.Close()

	h := BuildChecker([]ServiceConfig{
			&TCPServiceConfig{
				endpoint: "http://localhost:7777",
				name: "test-server0",
				timeout: defaultDuration,
			},
			&TCPServiceConfig{
				endpoint: "http://localhost:8888",
				name: "test-server1",
				timeout: defaultDuration,
			},
		},
		map[string]string{
			"version": "0.1.2",
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

	if hs.Status != ServiceStatusNOK {
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
