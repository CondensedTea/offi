package logstf

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var logsTfSearchTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:        "logstf_search_duration_seconds",
	Help:        "Response time of logs.tf API methods",
	ConstLabels: make(prometheus.Labels),
	Buckets: []float64{
		0.1, // 100 ms
		0.2,
		0.5,
		1.0, // 1s
		3.0,
		8.0,
		15.0,
		30.0,
	},
},
	[]string{},
)
