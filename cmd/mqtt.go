package main

import (
	"fmt"

	"github.com/wind-c/comqtt/v2/mqtt"
	"github.com/wind-c/comqtt/v2/mqtt/hooks/auth"
	"github.com/wind-c/comqtt/v2/mqtt/listeners"
)

func createBroker(cfg *config) (*mqtt.Server, error) {
	broker := mqtt.New(&mqtt.Options{
		InlineClient: true,
		Logger:       logger,
	})

	if cfg.MqttUser != "" && cfg.MqttPass != "" {
		broker.AddHook(new(auth.Hook), &auth.Options{
			Ledger: &auth.Ledger{
				Auth: auth.AuthRules{
					// {Remote: "127.0.0.1:*", Allow: true},
					// {Remote: "localhost:*", Allow: true},
					// {Remote: "::1:*", Allow: true},
					{Username: auth.RString(cfg.MqttHost), Password: auth.RString(cfg.MqttPass), Allow: true},
				},
				ACL: auth.ACLRules{},
			},
		})
	} else {
		if err := broker.AddHook(&auth.AllowHook{}, nil); err != nil {
			return nil, err
		}
	}

	address := fmt.Sprintf("%s:%d", cfg.MqttHost, cfg.MqttPort)
	tcpListener := listeners.NewTCP(fmt.Sprintf("mqtt-%d", cfg.MqttPort), address, nil)
	if err := broker.AddListener(tcpListener); err != nil {
		return nil, err
	}
	return broker, nil
}
