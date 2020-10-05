package health

import "testing"

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
	cfg := NewTCPChecker(
		"success-test",
		"http://localhost:9999",
		defaultDuration,
	)

	s := cfg.Test()

	t.Log(s.Details)
	if s.Status != ServiceStatusOK {
		t.Errorf("Status should be UP, but was '%s'", s.Status)
	}
}

func TestTcpTestHostUnavailable(t *testing.T) {
	cfg := NewTCPChecker(
		"invalid-test",
		"http://abc.xyz:1234",
		defaultDuration,
	)

	s := cfg.Test()
	if s.Status != ServiceStatusNOK {
		t.Errorf("Status deve estar DOWN, mas esta '%s'", s.Status)
	}
}

func TestTcpTestTargetInvalido(t *testing.T) {
	cfg := NewTCPChecker(
		"abc.com_invalid-test",
		"abc.com",
		defaultDuration,
	)
	s := cfg.Test()

	t.Log(s.Details)
	if s.Status != ServiceStatusNOK {
		t.Errorf("Status deve estar DOWN, mas esta '%s'", s.Status)
	}
}
