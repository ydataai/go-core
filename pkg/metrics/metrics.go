package metrics

import (
	"context"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

// Metrics includes metrics used in controller
type Metrics struct {
	client client.Client
	name   string
}

// NewClient register new service metrics
func NewClient(client client.Client, name string) *Metrics {
	return &Metrics{
		client: client,
		name:   name,
	}
}

// Creation is used when an object is created
func (m *Metrics) Creation() *prometheus.CounterVec {
	mc := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%v_create_total", m.name),
			Help: "Total times of creating",
		},
		[]string{"namespace"},
	)

	metrics.Registry.MustRegister(mc)
	return mc

}

// Running is used when an object is running
func (m *Metrics) Running() *prometheus.GaugeVec {
	mc := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("%v_running", m.name),
			Help: "Current running in the cluster",
		},
		[]string{"namespace"},
	)
	metrics.Registry.MustRegister(mc)
	return mc
}

// FailCreation is used when an object is failed
func (m *Metrics) FailCreation() *prometheus.CounterVec {
	mc := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%v_create_failed_total", m.name),
			Help: "Total failure times of creating",
		},
		[]string{"namespace"},
	)

	metrics.Registry.MustRegister(mc)
	return mc
}

// CullingCount is used when an object needs to be stopped
func (m *Metrics) CullingCount() *prometheus.CounterVec {
	mc := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: fmt.Sprintf("%v_culling_total", m.name),
			Help: "Total times of culling",
		},
		[]string{"namespace", "name"},
	)
	metrics.Registry.MustRegister(mc)
	return mc
}

// CullingTimestamp is used to set an annotation when the object is stopped
func (m *Metrics) CullingTimestamp() *prometheus.GaugeVec {
	mc := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: fmt.Sprintf("last_%v_culling_timestamp_seconds", m.name),
			Help: "Timestamp of the last culling in seconds",
		},
		[]string{"namespace", "name"},
	)
	metrics.Registry.MustRegister(mc)
	return mc
}

// Describe implements the prometheus.Collector interface.
func (m *Metrics) Describe(ch chan<- *prometheus.Desc) {
	m.Running().Describe(ch)
	m.Creation().Describe(ch)
	m.FailCreation().Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (m *Metrics) Collect(ch chan<- prometheus.Metric) {
	m.scrape()
	m.Running().Collect(ch)
	m.Creation().Collect(ch)
	m.FailCreation().Collect(ch)
}

func (m *Metrics) scrape() {
	stsList := &appsv1.StatefulSetList{}
	err := m.client.List(context.TODO(), stsList)
	if err != nil {
		return
	}
	stsCache := make(map[string]float64)
	for _, v := range stsList.Items {
		name, ok := v.Spec.Template.GetLabels()[m.name+"-name"]
		if ok && name == v.Name {
			stsCache[v.Namespace]++
		}
	}

	for ns, v := range stsCache {
		m.Running().WithLabelValues(ns).Set(v)
	}
}
