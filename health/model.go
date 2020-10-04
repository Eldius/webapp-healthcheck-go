package health

type ServiceStatus string

const (
	ServiceStatusOK  ServiceStatus = "UP"
	ServiceStatusNOK ServiceStatus = "DOWN"
)

type Status struct {
	Name    string
	Status  ServiceStatus
	Details map[string]string
}

type HealthStatus struct {
	Version   string
	BuildDate string
	Statuses  []Status
}

func (hs *HealthStatus) Status() ServiceStatus {
	for _, s := range hs.Statuses {
		if s.Status != ServiceStatusOK {
			return ServiceStatusNOK
		}
	}
	return ServiceStatusOK
}
