package health

import (
	"testing"
	"time"
	"gopkg.in/h2non/gock.v1"
)

func TestHttpCheckOk(t *testing.T)  {
	defer gock.Off()
	gock.New("http://test-host.xpto").
		Get("/foo/123").
		Reply(200).
		JSON(map[string]string{"foo": "bar"})

	c := NewHTTPChecker("http-check", "http://test-host.xpto/foo/123", time.Duration(1 * time.Second), 200)

	s := c.Test()

	if s.Status != CheckerStatusOK {
		t.Errorf("Should return status 'UP', bu returned '%s'", s.Status)
	}
}

func TestHttpCheckConnectivityFailure(t *testing.T)  {
	c := NewHTTPChecker("http-check", "http://test-host.xpto/foo/123", time.Duration(1 * time.Second), 200)

	s := c.Test()

	t.Log(s.Details)
	if s.Status != CheckerStatusNOK {
		t.Errorf("Should return status 'DOWN', bu returned '%s'", s.Status)
	}
}
