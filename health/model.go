package health

/*
ServiceStatus is the result of health check
possible values are 'UP' and 'DOWN'
*/
type ServiceStatus string

const (
	// ServiceStatusOK service is OK
	ServiceStatusOK  ServiceStatus = "UP"
	// ServiceStatusNOK service is not OK
	ServiceStatusNOK ServiceStatus = "DOWN"
)

/*
Status is the check result for a service
*/
type Status struct {
	Name    string
	Status  ServiceStatus
	Details map[string]string
}

/*
HealthStatus is the healthcheck status
aggregate all service status
*/
type HealthStatus struct {
	Info map[string]string
	Statuses  []Status
}

/*
Status checks if all services are ok
*/
func (hs *HealthStatus) Status() ServiceStatus {
	for _, s := range hs.Statuses {
		if s.Status != ServiceStatusOK {
			return ServiceStatusNOK
		}
	}
	return ServiceStatusOK
}
