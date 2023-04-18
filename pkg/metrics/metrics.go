package metrics

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cloud-native-observability/metrics-simulator/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
)

var generator *Generator

type Generator struct {
	config     *config.GeneratorConfig
	registry   *prometheus.Registry
	metrics    []MetricGenerator
	metricsMap map[string]MetricGenerator
}
type MetricGenerator interface {
	generate()
}
type CounterMetric struct {
	config *config.CounterConfig
	metric prometheus.Counter
}

func (m *CounterMetric) generate() {
	go func() {
		for {
			m.metric.Add(float64(rand.Intn(10)))
			time.Sleep(time.Second) // FIXME: use config value
		}
	}()
}

type GaugeMetric struct {
	config *config.GaugeConfig
	metric prometheus.Gauge
	value  float64
}

func (m *GaugeMetric) generate() {
	go func() {
		for {
			delta := float64(rand.Intn(20) - 10)
			if m.value+delta > float64(m.config.Range.Upper) || m.value+delta < float64(m.config.Range.Lower) {
				delta = 0.0
			}
			m.value = m.value + delta
			m.metric.Set(m.value)
			time.Sleep(time.Second) // FIXME: use config value
		}
	}()
}
func NewGenerator(file string, registry *prometheus.Registry) (*Generator, error) {
	config, err := config.ParseConfig(file)
	generator = &Generator{config: config,
		registry:   registry,
		metrics:    []MetricGenerator{},
		metricsMap: make(map[string]MetricGenerator),
	}
	if err != nil {
		return nil, err
	}
	return generator, nil
}
func (g *Generator) Load() {
	g.createMetricsGenerator()
	g.generateMetrics()
}

func (g *Generator) generateMapkey(mname string, labels map[string]string) string {
	key := mname

	for k, v := range labels {
		key = key + ":" + k + "," + v

	}
	return key

}
func (g *Generator) createMetricsGenerator() {
	if g.config.Counters != nil {
		for _, counter := range g.config.Counters {
			for i := 0; i < counter.Number; i++ {

				labels := make(map[string]string)
				for _, lsetting := range counter.Labels {
					v := selectRandomLabelValue(lsetting.ValueSet)
					labels[lsetting.Name] = v
				}
				opts := prometheus.CounterOpts{}
				if counter.Number == 1 {
					opts.Name = counter.Prefix
				} else {
					opts.Name = fmt.Sprint(counter.Prefix, "_", i)
				}
				opts.ConstLabels = labels
				pmetric := prometheus.NewCounter(opts)
				cmetric := CounterMetric{
					config: counter,
					metric: pmetric,
				}
				g.registry.MustRegister(pmetric)
				g.metrics = append(g.metrics, &cmetric)
				g.metricsMap[g.generateMapkey(opts.Name, opts.ConstLabels)] = &cmetric
			}
		}
	}
	if g.config.Gauges != nil {
		for _, gauge := range g.config.Gauges {
			for i := 0; i < gauge.Number; i++ {
				labels := make(map[string]string)
				for _, lsetting := range gauge.Labels {
					v := selectRandomLabelValue(lsetting.ValueSet)
					labels[lsetting.Name] = v
				}
				opts := prometheus.GaugeOpts{}
				opts.Name = fmt.Sprint(gauge.Prefix, "_", i)
				opts.ConstLabels = labels
				pmetric := prometheus.NewGauge(opts)
				g.registry.MustRegister(pmetric)
				cmetric := GaugeMetric{config: gauge,
					metric: pmetric,
					value:  float64(gauge.Range.Lower),
				}
				g.metrics = append(g.metrics, &cmetric)
				g.metricsMap[g.generateMapkey(opts.Name, opts.ConstLabels)] = &cmetric
			}
		}
	}
}

func (g *Generator) generateMetrics() {
	for _, m := range g.metrics {
		m.generate()
	}
}
func selectRandomLabelValue(labelValues []string) string {
	if len(labelValues) == 0 {
		return ""
	}
	i := rand.Intn(len(labelValues))
	return labelValues[i]
}
