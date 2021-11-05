package k8s

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

// Metrics is a simple struct with some metric templates
type Metrics struct {
	client             client.Client
	Create             *prometheus.CounterVec
	FailedCreate       *prometheus.CounterVec
	Running            *prometheus.GaugeVec
	CullingCount       *prometheus.CounterVec
	FailedCullingCount *prometheus.CounterVec
	CullingTimestamp   *prometheus.GaugeVec
}

// NewMetrics register new metrics service
// WARNING: you cannot call this method twice passing the same name
// otherwise, it can cause panic
func NewMetrics(client client.Client, name string) *Metrics {
	mc := &Metrics{
		client: client,
		Create: newCounter(
			fmt.Sprintf("%v_create_total", name),
			[]string{"namespace", "kind"},
		),
		FailedCreate: newCounter(
			fmt.Sprintf("%v_create_failed_total", name),
			[]string{"namespace", "kind"},
		),
		Running: newGauge(
			fmt.Sprintf("%v_running", name),
			[]string{"namespace", "kind"},
		),
		CullingCount: newCounter(
			fmt.Sprintf("%v_culling_total", name),
			[]string{"namespace", "kind"},
		),
		FailedCullingCount: newCounter(
			fmt.Sprintf("%v_culling_failed_total", name),
			[]string{"namespace", "kind"},
		),
		CullingTimestamp: newGauge(
			fmt.Sprintf("last_%v_culling_timestamp_seconds", name),
			[]string{"namespace", "kind"},
		),
	}

	metrics.Registry.MustRegister(mc)
	return mc
}

// CustomCount now you can create your own gauge metric
// WARNING: you cannot call this method twice passing the same metric name
// otherwise, it can cause panic
func (m *Metrics) CustomCount(name string, labels []string) *prometheus.CounterVec {
	cc := newCounter(name, labels)

	metrics.Registry.MustRegister(cc)
	return cc
}

// CustomGauge now you can create your own gauge metric
// WARNING: you cannot call this method twice passing the same metric name
// otherwise, it can cause panic
func (m *Metrics) CustomGauge(name string, labels []string) *prometheus.GaugeVec {
	cc := newGauge(name, labels)

	metrics.Registry.MustRegister(cc)
	return cc
}

// Describe implements the prometheus.Collector interface.
func (m *Metrics) Describe(ch chan<- *prometheus.Desc) {
	m.Running.Describe(ch)
	m.Create.Describe(ch)
	m.FailedCreate.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (m *Metrics) Collect(ch chan<- prometheus.Metric) {
	m.Running.Collect(ch)
	m.Create.Collect(ch)
	m.FailedCreate.Collect(ch)
}
