package health

import "testing"

func TestHealthStatusOK(t *testing.T) {
	h := HealthStatus{
		Services: []Status{
			{
				Name:   "service01",
				Status: CheckerStatusOK,
			},
			{
				Name:   "service02",
				Status: CheckerStatusOK,
			},
			{
				Name:   "service03",
				Status: CheckerStatusOK,
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
				Status: CheckerStatusNOK,
			},
			{
				Name:   "service02",
				Status: CheckerStatusOK,
			},
			{
				Name:   "service03",
				Status: CheckerStatusOK,
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
				Status: CheckerStatusOK,
			},
			{
				Name:   "service02",
				Status: CheckerStatusNOK,
			},
			{
				Name:   "service03",
				Status: CheckerStatusOK,
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
				Status: CheckerStatusOK,
			},
			{
				Name:   "service02",
				Status: CheckerStatusOK,
			},
			{
				Name:   "service03",
				Status: CheckerStatusNOK,
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
				Status: CheckerStatusNOK,
			},
			{
				Name:   "service02",
				Status: CheckerStatusNOK,
			},
			{
				Name:   "service03",
				Status: CheckerStatusNOK,
			},
		},
	}

	h.Evaluate()

	if h.Status != "DOWN" {
		t.Errorf("h.Status should return 'DOWN', but returned '%s'", h.Status)
	}
}
