package client

// Application contains list of the Applications associated with your New Relic account
type Application struct {
	AppsList []appsList `json:"applications"`
}

// appsList the list of application IDs or the application language as reported by the agents
type appsList struct {
	AppID      int    `json:"id"`
	AppName    string `json:"name"`
	AppLang    string `json:"language"`
	AppHealth  string `json:"health_status"`
	AppSummary appSummary
}

// appSummary basic information of app
type appSummary struct {
	ResponseTime  float32 `json:"response_time"`
	Throughput    float32 `json:"throughput"`
	ErrorRate     float32 `json:"error_rate"`
	ApdexTarget   float32 `json:"apdex_target"`
	ApdexScore    float32 `json:"apdex_score"`
	HostCount     float32 `json:"host_count"`
	InstanceCount float32 `json:"instance_count"`
}

// ApplicationHost application hosts can be filtered by hostname, or the list of application host IDs.
type ApplicationHost struct {
	AppsHosts []appHost `json:"application_hosts"`
}

// appHost basic information application host
type appHost struct {
	HostID      int        `json:"id"`
	AppName     string     `json:"application_name"`
	Host        string     `json:"host"`
	HostSummary appSummary `json:"application_summary"`
}

// MetricsData the list of available metrics can be returned using the Metric Name API endpoint.
type MetricsData struct {
	Data metricsData `json:"metric_data"`
}

// metricsData collection by time
type metricsData struct {
	From    string    `json:"from"`
	To      string    `json:"to"`
	Metrics []metrics `json:"metrics"`
}

// metrics by name and for specific time
type metrics struct {
	Name       string              `json:"name"`
	Timeslices []metricsTimeslices `json:"timeslices"`
}

// metricsTimeslices data collected in time
type metricsTimeslices struct {
	Values timeslicesValue `json:"values"`
}

// timeslicesValue metrics data in time
type timeslicesValue struct {
	Requests float32 `json:"requests_per_minute"`
	Errors   float32 `json:"errors_per_minute"`
}
