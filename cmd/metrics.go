package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	commonLabelKeys = []string{"device", "index"}
)

func labelsFromConfig(cfg *config) prometheus.Labels {
	labels := prometheus.Labels{}
	for _, label := range cfg.PrometheusLabels {
		parts := strings.SplitN(label, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" || value == "" {
			continue
		}
		labels[key] = value
	}
	return labels
}

type prometheusCollector struct {
	registry            *prometheus.Registry
	constLabels         prometheus.Labels
	OutputGauge         *prometheus.GaugeVec
	APowerGauge         *prometheus.GaugeVec
	VoltageGauge        *prometheus.GaugeVec
	CurrentGauge        *prometheus.GaugeVec
	FreqGauge           *prometheus.GaugeVec
	EnergyTotalGauge    *prometheus.GaugeVec
	RetEnergyTotalGauge *prometheus.GaugeVec
	EnergyMin0Gauge     *prometheus.GaugeVec
	RetEnergyMin0Gauge  *prometheus.GaugeVec
	TempGauge           *prometheus.GaugeVec
}

func newPrometheusCollector(cfg *config) *prometheusCollector {
	labels := labelsFromConfig(cfg)
	c := &prometheusCollector{
		registry:            prometheus.NewRegistry(),
		constLabels:         labels,
		OutputGauge:         prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "shelly_output", Help: "Switch output state"}, commonLabelKeys),
		APowerGauge:         prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "shelly_apower", Help: "Active power in Watts"}, commonLabelKeys),
		VoltageGauge:        prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "shelly_voltage", Help: "Voltage in Volts"}, commonLabelKeys),
		CurrentGauge:        prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "shelly_current", Help: "Current in Amps"}, commonLabelKeys),
		FreqGauge:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "shelly_freq", Help: "Frequency in Hz"}, commonLabelKeys),
		EnergyTotalGauge:    prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "shelly_energy_total_wh", Help: "Total energy in Wh"}, commonLabelKeys),
		RetEnergyTotalGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "shelly_returned_energy_total_wh", Help: "Returned energy in Wh"}, commonLabelKeys),
		EnergyMin0Gauge:     prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "shelly_energy_minute_0_mwh", Help: "Energy in mWh from minute -1"}, commonLabelKeys),
		RetEnergyMin0Gauge:  prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "shelly_returned_energy_minute_0_mwh", Help: "Returned energy in mWh from minute -1"}, commonLabelKeys),
		TempGauge:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "shelly_temperature_celsius", Help: "Temperature in Celsius"}, commonLabelKeys),
	}
	c.registry.MustRegister(
		c.OutputGauge,
		c.APowerGauge,
		c.VoltageGauge,
		c.CurrentGauge,
		c.FreqGauge,
		c.EnergyTotalGauge,
		c.RetEnergyTotalGauge,
		c.EnergyMin0Gauge,
		c.RetEnergyMin0Gauge,
		c.TempGauge,
	)
	return c
}

func (c *prometheusCollector) Registry() *prometheus.Registry {
	return c.registry
}

func (c *prometheusCollector) Collect(payload []byte, labels []string) error {
	var d ShellyData
	json.Unmarshal(payload, &d)

	if len(commonLabelKeys) != len(labels) {
		return fmt.Errorf("expected %d labels, got %d", len(commonLabelKeys), len(labels))
	}

	l := prometheus.Labels{
		commonLabelKeys[0]: labels[0],
		commonLabelKeys[1]: labels[1],
	}

	outputGauge := c.OutputGauge.With(l)
	if d.Output {
		outputGauge.Set(1)
	} else {
		outputGauge.Set(0)
	}
	if d.ActivePower != nil {
		c.APowerGauge.With(l).Set(*d.ActivePower)
	}
	if d.Voltage != nil {
		c.VoltageGauge.With(l).Set(*d.Voltage)
	}
	if d.Current != nil {
		c.CurrentGauge.With(l).Set(*d.Current)
	}
	if d.Freq != nil {
		c.FreqGauge.With(l).Set(*d.Freq)
	}
	if d.ActiveEnergy != nil {
		c.EnergyTotalGauge.With(l).Set(d.ActiveEnergy.Total)
		if len(d.ActiveEnergy.ByMinute) > 0 {
			c.EnergyMin0Gauge.With(l).Set(d.ActiveEnergy.ByMinute[0])
		}
	}
	if d.ReturnedActiveEnergy != nil {
		c.RetEnergyTotalGauge.With(l).Set(d.ReturnedActiveEnergy.Total)
		if len(d.ReturnedActiveEnergy.ByMinute) > 0 {
			c.RetEnergyMin0Gauge.With(l).Set(d.ReturnedActiveEnergy.ByMinute[0])
		}
	}
	if d.Temperature != nil {
		c.TempGauge.With(l).Set(d.Temperature.TC)
	}

	return nil
}
