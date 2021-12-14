package metrics

import "github.com/prometheus/client_golang/prometheus"

var EmailHits = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "email_hits",
}, []string{"status", "operation"})
