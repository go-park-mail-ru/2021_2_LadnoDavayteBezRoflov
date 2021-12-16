package metrics

import "github.com/prometheus/client_golang/prometheus"

var SessionHits = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "session_hits",
}, []string{"status", "operation"})
