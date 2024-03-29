package health

import "testing"

func TestStatusOK(t *testing.T) {
	h := Status{
		Services: []ServiceStatus{
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

func TestStatusNOK0(t *testing.T) {
	h := Status{
		Services: []ServiceStatus{
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

func TestStatusNOK1(t *testing.T) {
	h := Status{
		Services: []ServiceStatus{
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

func TestStatusNOK2(t *testing.T) {
	h := Status{
		Services: []ServiceStatus{
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

func TestStatusNOK3(t *testing.T) {
	h := Status{
		Services: []ServiceStatus{
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
