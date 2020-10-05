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
	Status    ServiceStatus
	Info      map[string]string
	Services  []Status
}

/*
Evaluate checks all services status and set health status
*/
func (hs *HealthStatus) Evaluate() {
	hs.Status = ServiceStatusOK
	for _, s := range hs.Services {
		if s.Status != ServiceStatusOK {
			hs.Status = ServiceStatusNOK
		}
	}
}
