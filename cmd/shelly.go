package main

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/wind-c/comqtt/v2/mqtt"
	mqttPackets "github.com/wind-c/comqtt/v2/mqtt/packets"
)

/** Example payload:

```
{
    "id": 0,
    "source": "WS_in",
    "output": false,
    "apower": 0.0,
    "voltage": 241.4,
    "freq": 50.0,
    "current": 0.000,
    "aenergy": {
        "total": 161.961,
        "by_minute": [
            0.000,
            0.000,
            0.000
        ],
        "minute_ts": 1744214220
    },
    "ret_aenergy": {
        "total": 0.000,
        "by_minute": [
            0.000,
            0.000,
            0.000
        ],
        "minute_ts": 1744214220
    },
    "temperature": {
        "tC": 31.5,
        "tF": 88.8
    }
}
```
*/

type ShellyDataEnergyCounter struct {
	// Total energy consumed in Watt-hours
	Total float64 `json:"total"`
	// Total energy flow in Milliwatt-hours for the last three complete minutes. The 0-th element indicates the counts accumulated during the minute preceding minute_ts.
	ByMinute []float64 `json:"by_minute"`
	// Unix timestamp marking the start of the current minute (in UTC).
	MinuteTs int64 `json:"minute_ts"`
}

type ShellyDataTemperature struct {
	// Temperature in Celsius
	TC float64 `json:"tC"`
	// Temperature in Fahrenheit
	TF float64 `json:"tF"`
}

type ShellyData struct {
	// Id of the Switch component instance
	ID int `json:"id"`
	// Source of the last command, for example: init, WS_in, http, ...
	Source string `json:"source"`
	// true if the output channel is currently on, false otherwise
	Output bool `json:"output"`
	// Unix timestamp, start time of the timer (in UTC) (shown if the timer is triggered)
	TimerStartedAt *int64 `json:"timer_started_at"`
	// Duration of the timer in seconds (shown if the timer is triggered)
	TimerDuration *int `json:"timer_duration"`
	// Last measured instantaneous active power (in Watts) delivered to the attached load (shown if applicable)
	ActivePower *float64 `json:"apower"`
	// Last measured voltage in Volts (shown if applicable)
	Voltage *float64 `json:"voltage"`
	// Last measured current in Amperes (shown if applicable)
	Current *float64 `json:"current"`
	// Last measured power factor (shown if applicable)
	PF *float64 `json:"pf"`
	// Last measured frequency in Hz (shown if applicable)
	Freq *float64 `json:"freq"`
	// Information about the active energy counter (shown if applicable)
	ActiveEnergy *ShellyDataEnergyCounter `json:"aenergy"`
	// Information about the returned active energy counter * (shown if applicable)
	ReturnedActiveEnergy *ShellyDataEnergyCounter `json:"ret_aenergy"`
	// Information about the temperature (shown if applicable)
	Temperature *ShellyDataTemperature `json:"temperature"`
	// Error conditions occurred. May contain overtemp, overpower, overvoltage, undervoltage, (shown if at least one error is present)
	Errors *[]string `json:"errors"`
}

func shellyHandler(cl *mqtt.Client, sub mqttPackets.Subscription, pk mqttPackets.Packet) {
	if strings.Index(pk.TopicName, "$SYS") == 0 {
		// ignore system topics
		return
	}
	// parse host and device index from topic name
	regex := regexp.MustCompile(`^([^/]+)/status/switch:(\d+)$`)
	matches := regex.FindStringSubmatch(pk.TopicName)
	if len(matches) < 3 {
		// unknown topic
		return
	} else {
		device := matches[1]
		index := matches[2]
		var data ShellyData
		json.Unmarshal(pk.Payload, &data)
		if err := collector.Collect(data, []string{device, index}); err != nil {
			logger.Error("failed to collect metrics", "device", device, "index", index, "error", err)
		}
	}
}
