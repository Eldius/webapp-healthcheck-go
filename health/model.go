package health

/*
CheckerStatus is the result of health check
possible values are 'UP' and 'DOWN'
*/
type CheckerStatus string

const (
	// CheckerStatusOK service is OK
	CheckerStatusOK CheckerStatus = "UP"
	// CheckerStatusNOK service is not OK
	CheckerStatusNOK CheckerStatus = "DOWN"
)

/*
Status is the check result for a service
*/
type Status struct {
	Name    string `json:"name"`
	Status  CheckerStatus `json:"status"`
	Details map[string]string `json:"details"`
}

/*
HealthStatus is the healthcheck status
aggregate all service status
*/
type HealthStatus struct {
	Status   CheckerStatus `json:"status"`
	Info     map[string]string `json:"info"`
	Services []Status `json:"services"`
}

/*
Evaluate checks all services status and set health status
*/
func (hs *HealthStatus) Evaluate() {
	hs.Status = CheckerStatusOK
	for _, s := range hs.Services {
		if s.Status != CheckerStatusOK {
			hs.Status = CheckerStatusNOK
		}
	}
}
