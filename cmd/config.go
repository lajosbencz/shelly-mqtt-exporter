package main

type config struct {
	MqttHost         string   `env:"MQTT_HOST" envDefault:"localhost"`
	MqttPort         int      `env:"MQTT_PORT" envDefault:"1883"`
	MqttUser         string   `env:"MQTT_USER" envDefault:""`
	MqttPass         string   `env:"MQTT_PASS" envDefault:""`
	PrometheusHost   string   `env:"PROM_HOST" envDefault:"localhost"`
	PrometheusPort   int      `env:"PROM_PORT" envDefault:"2112"`
	PrometheusPath   string   `env:"PROM_PATH" envDefault:"/metrics"`
	PrometheusUser   string   `env:"PROM_USER" envDefault:""`
	PrometheusPass   string   `env:"PROM_PASS" envDefault:""`
	PrometheusLabels []string `env:"PROM_LABELS" envDefault:""`
}
