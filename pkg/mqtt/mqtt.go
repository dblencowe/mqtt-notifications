package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
)

type BrokerOption func(*BrokerConfig) error

type BrokerConfig struct {
	useAuth  bool
	username string
	password string
}

func OptionAuth(username, password string) BrokerOption {
	return func(cfg *BrokerConfig) error {
		if len(username) != 0 || len(password) != 0 {
			cfg.useAuth = true
			cfg.username = username
			cfg.password = password
		}
		return nil
	}
}

var messageQueue chan mqtt.Message

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	opts := client.OptionsReader()
	log.Printf("connected to %s as %s", opts.Servers()[0].Host, opts.ClientID())
}

var disconnectHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	opts := client.OptionsReader()
	log.Printf("disconnected from %s: %s", opts.Servers()[0].Host, err)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	messageQueue <- msg
}

func Connect(uri string, port int, cfgOpts ...BrokerOption) (mqtt.Client, error) {
	messageQueue = make(chan mqtt.Message)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", uri, port))
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	opts.SetClientID(fmt.Sprintf("mqtt-notifications-%s", name))
	cfg := &BrokerConfig{}
	for _, op := range cfgOpts {
		err := op(cfg)
		if err != nil {
			return nil, err
		}
	}
	if cfg.useAuth {
		opts.SetUsername(cfg.username)
		opts.SetPassword(cfg.password)
	}
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = disconnectHandler
	opts.SetDefaultPublishHandler(messagePubHandler)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return client, nil
}

func Publish(client mqtt.Client, subject, value string) error {
	name, err := os.Hostname()
	if err != nil {
		return err
	}

	if token := client.Publish(fmt.Sprintf("mqtt-notifications/%s/%s", name, subject), 0, false, value); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func Subscribe(client mqtt.Client, topic string) error {
	if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Printf("subscribed to topic %s", topic)
	return nil
}

func Receive() chan mqtt.Message {
	return messageQueue
}
