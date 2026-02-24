package connection

import (
	"encore.app/internal/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func NewMQTT(cfg *config.Config) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.MqttBrokker)
	opts.SetClientID(cfg.MqttClientId)
	opts.SetCleanSession(true)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}
