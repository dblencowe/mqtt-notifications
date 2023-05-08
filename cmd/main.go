package main

import (
	"encoding/json"
	"git.strawberryelk.internal/dblencowe/mqtt-notifications/pkg/mqtt"
	"github.com/caarlos0/env/v8"
	"github.com/gen2brain/beeep"
	"log"
)

type config struct {
	BrokerUrl      string `env:"MQTT_BROKER_URL,required"`
	BrokerPort     int    `env:"MQTT_BROKER_PORT,required"`
	BrokerUsername string `env:"MQTT_BROKER_USERNAME"`
	BrokerPassword string `env:"MQTT_BROKER_PASSWORD"`
	BrokerTopic    string `env:"MQTT_NOTIFICATIONS_TOPIC,required"`
}

type notificationMessage struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("unable to parse environment config ", err)
	}

	var mqttOpts []mqtt.BrokerOption
	if len(cfg.BrokerUsername) > 0 || len(cfg.BrokerPassword) > 0 {
		mqttOpts = append(mqttOpts, mqtt.OptionAuth(cfg.BrokerUsername, cfg.BrokerPassword))
	}
	client, err := mqtt.Connect(cfg.BrokerUrl, cfg.BrokerPort, mqttOpts...)
	if err != nil {
		log.Fatal("cannot connect to mqtt server ", err)
	}
	defer client.Disconnect(0)

	err = mqtt.Subscribe(client, cfg.BrokerTopic)
	if err != nil {
		log.Fatal("cannot subscribe to mqtt topic ", err)
	}
	for msg := range mqtt.Receive() {
		var notif notificationMessage
		err := json.Unmarshal(msg.Payload(), &notif)
		if err != nil {
			log.Printf("error unmarshalling message: %s", err)
			continue
		}
		err = beeep.Notify(notif.Title, notif.Message, "")
		if err != nil {
			log.Printf("error displaying notification: %s", err)
		}
	}
}
