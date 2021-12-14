package metrics

import "github.com/penglongli/gin-metrics/ginmetrics"

var APIErrors = &ginmetrics.Metric{
	Type:        ginmetrics.Counter,
	Name:        "api_errors",
	Description: "api errors",
	Labels:      []string{"status", "description"},
}
