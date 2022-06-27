package metrics

import "github.com/prometheus/client_golang/prometheus"

type GaugeMap map[string]*LocationGauges

type LocationGauges struct {
	Temperature    prometheus.Gauge
	MinTemperature prometheus.Gauge
	MaxTemperature prometheus.Gauge
	Humidity       prometheus.Gauge
	Pressure       prometheus.Gauge
}

func NewLocationGauges(locID string) *LocationGauges {
	gauges := &LocationGauges{}
	labels := map[string]string{"id": locID}
	type opts struct {
		dest *prometheus.Gauge
		name string
	}
	for _, opt := range []opts{
		{&gauges.Temperature, "temperature"},
		{&gauges.MinTemperature, "min_temperature"},
		{&gauges.MaxTemperature, "max_temperature"},
		{&gauges.Humidity, "humidity"},
		{&gauges.Pressure, "pressure"},
	} {
		*opt.dest = prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        opt.name,
			ConstLabels: labels,
		})
	}
	return gauges
}

func (l *LocationGauges) MustRegister(reg prometheus.Registerer) {
	reg.MustRegister(l.Temperature)
	reg.MustRegister(l.MinTemperature)
	reg.MustRegister(l.MaxTemperature)
	reg.MustRegister(l.Humidity)
	reg.MustRegister(l.Pressure)
}
