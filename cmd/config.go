package main

import (
	"crypto/sha256"

	"github.com/caarlos0/env"
)

func getConfig() (*config, error) {
	cfg := &config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	if cfg.PrometheusPass != "" {
		sha := sha256.Sum256([]byte(cfg.PrometheusPass))
		cfg.PrometheusPassSha = &sha
		cfg.PrometheusPass = ""
	}
	if len(cfg.PrometheusPath) > 0 && cfg.PrometheusPath[0] != '/' {
		cfg.PrometheusPath = "/" + cfg.PrometheusPath
	}
	return cfg, nil
}

type config struct {
	MqttHost              string `env:"MQTT_HOST" envDefault:"0.0.0.0"`
	MqttPort              int    `env:"MQTT_PORT" envDefault:"1883"`
	MqttUser              string `env:"MQTT_USER" envDefault:""`
	MqttPass              string `env:"MQTT_PASS" envDefault:""`
	PrometheusHost        string `env:"PROM_HOST" envDefault:"0.0.0.0"`
	PrometheusPort        int    `env:"PROM_PORT" envDefault:"2112"`
	PrometheusPath        string `env:"PROM_PATH" envDefault:"/metrics"`
	PrometheusUser        string `env:"PROM_USER" envDefault:""`
	PrometheusPass        string `env:"PROM_PASS" envDefault:""`
	PrometheusPassSha     *[32]byte
	PrometheusLabels      string `env:"PROM_LABELS" envDefault:""`
	PrometheusTLSCertPath string `env:"HTTP_TLS_CERT" envDefault:""`
	PrometheusTLSKeyPath  string `env:"HTTP_TLS_KEY" envDefault:""`
	MqttTLSCertPath       string `env:"MQTT_TLS_CERT" envDefault:""`
	MqttTLSKeyPath        string `env:"MQTT_TLS_KEY" envDefault:""`
}

func (c *config) IsMqttTls() bool {
	return c.MqttTLSCertPath != "" && c.MqttTLSKeyPath != ""
}

func (c *config) IsPromTls() bool {
	return c.PrometheusTLSCertPath != "" && c.PrometheusTLSKeyPath != ""
}
