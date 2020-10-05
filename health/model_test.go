package health

import "testing"

func TestHealthStatusOK(t *testing.T) {
	h := HealthStatus{
		Services: []Status{
			{
				Name:   "service01",
				Status: ServiceStatusOK,
			},
			{
				Name:   "service02",
				Status: ServiceStatusOK,
			},
			{
				Name:   "service03",
				Status: ServiceStatusOK,
			},
		},
	}

	h.Evaluate()

	if h.Status != "UP" {
		t.Errorf("h.Status should return 'UP', but returned '%s'", h.Status)
	}
}

func TestHealthStatusNOK0(t *testing.T) {
	h := HealthStatus{
		Services: []Status{
			{
				Name:   "service01",
				Status: ServiceStatusNOK,
			},
			{
				Name:   "service02",
				Status: ServiceStatusOK,
			},
			{
				Name:   "service03",
				Status: ServiceStatusOK,
			},
		},
	}

	h.Evaluate()

	if h.Status != "DOWN" {
		t.Errorf("h.Status should return 'DOWN', but returned '%s'", h.Status)
	}
}

func TestHealthStatusNOK1(t *testing.T) {
	h := HealthStatus{
		Services: []Status{
			{
				Name:   "service01",
				Status: ServiceStatusOK,
			},
			{
				Name:   "service02",
				Status: ServiceStatusNOK,
			},
			{
				Name:   "service03",
				Status: ServiceStatusOK,
			},
		},
	}

	h.Evaluate()

	if h.Status != "DOWN" {
		t.Errorf("h.Status should return 'DOWN', but returned '%s'", h.Status)
	}
}

func TestHealthStatusNOK2(t *testing.T) {
	h := HealthStatus{
		Services: []Status{
			{
				Name:   "service01",
				Status: ServiceStatusOK,
			},
			{
				Name:   "service02",
				Status: ServiceStatusOK,
			},
			{
				Name:   "service03",
				Status: ServiceStatusNOK,
			},
		},
	}

	h.Evaluate()

	if h.Status != "DOWN" {
		t.Errorf("h.Status should return 'DOWN', but returned '%s'", h.Status)
	}
}

func TestHealthStatusNOK3(t *testing.T) {
	h := HealthStatus{
		Services: []Status{
			{
				Name:   "service01",
				Status: ServiceStatusNOK,
			},
			{
				Name:   "service02",
				Status: ServiceStatusNOK,
			},
			{
				Name:   "service03",
				Status: ServiceStatusNOK,
			},
		},
	}

	h.Evaluate()

	if h.Status != "DOWN" {
		t.Errorf("h.Status should return 'DOWN', but returned '%s'", h.Status)
	}
}
