# Shelly MQTT Exporter

Shelly MQTT Exporter is a lightweight tool that collects data from Shelly devices via MQTT and exposes it as Prometheus metrics. This allows you to monitor and analyze your Shelly devices' performance and energy usage in a Prometheus-compatible monitoring system.


## Features

- Collects data from Shelly devices using MQTT.
- Exposes metrics such as power consumption, voltage, current, frequency, temperature, and more.
- Supports Prometheus scraping with optional basic authentication.
- Configurable via environment variables.


## Installation

Check the [releases page on Github](https://github.com/lajosbencz/shelly-mqtt-exporter/releases)


## Configuration

The exporter is configured using environment variables. Below are the available options:

| Environment Variable      | Default Value   | Description                                      |
|---------------------------|-----------------|--------------------------------------------------|
| `MQTT_HOST`               | `localhost`     | MQTT broker host.                               |
| `MQTT_PORT`               | `1883`          | MQTT broker port.                               |
| `MQTT_USER`               |                 | MQTT username (optional).                       |
| `MQTT_PASS`               |                 | MQTT password (optional).                       |
| `PROM_HOST`               | `localhost`     | Prometheus metrics server host.                 |
| `PROM_PORT`               | `2112`          | Prometheus metrics server port.                 |
| `PROM_PATH`               | `/metrics`      | Path for Prometheus metrics.                    |
| `PROM_USER`               |                 | Prometheus basic auth (optional).               |
| `PROM_PASS`               |                 | Prometheus basic auth password (optional).      |
| `PROM_LABELS`             |                 | Comma separated labels, for eg: foo=bar,baz=bax (optional).  |
| `LOG_LEVEL`               | `WARN`          | Logging level (`DEBUG`, `INFO`, `WARN`, `ERROR`).|
| `TLS_CERT`                |                 | Path to TLS certificate file (optional).        |
| `TLS_KEY`                 |                 | Path to TLS key file (optional).                |


## Metrics

The following metrics are exposed by the exporter:

- `shelly_output`: Switch output state (0 or 1).
- `shelly_apower`: Active power in Watts.
- `shelly_voltage`: Voltage in Volts.
- `shelly_current`: Current in Amperes.
- `shelly_freq`: Frequency in Hz.
- `shelly_energy_total_wh`: Total energy consumed in Watt-hours.
- `shelly_returned_energy_total_wh`: Total returned energy in Watt-hours.
- `shelly_energy_minute_0_mwh`: Energy in milliwatt-hours from the last minute.
- `shelly_returned_energy_minute_0_mwh`: Returned energy in milliwatt-hours from the last minute.
- `shelly_temperature_celsius`: Temperature in Celsius.


## Usage

1. Start the MQTT broker and ensure your Shelly devices are publishing data to it.
2. Run the exporter with the appropriate environment variables set.
3. Configure Prometheus to scrape metrics from the exporter.

Example Prometheus scrape configuration:
```yaml
scrape_configs:
  - job_name: "shelly-mqtt-exporter"
    static_configs:
      - targets: ["localhost:2112"]
```

## TLS

If you want to use TLS, you need to set the `TLS_CERT` and `TLS_KEY` environment variables to the paths of your TLS certificate and key files, respectively. Both MQTT and Prometheus connections will use these certificates for TLS.


## License

This project is licensed under the [MIT License](LICENSE).
