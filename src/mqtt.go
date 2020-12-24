package ttrn

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// Mqtt bundles the MQTT related stuff.
type Mqtt struct {
	broker        mqtt.Client
	turnoutEvents chan TurnoutEvent
	trainEvents   chan TrainEvent
}

// NewMqtt returns a new Mqtt instance.
func NewMqtt(cfg Config, turnoutEvents chan TurnoutEvent, trainEvents chan TrainEvent) Mqtt {
	rsl := Mqtt{
		turnoutEvents: turnoutEvents,
		trainEvents:   trainEvents,
	}

	opt := mqtt.NewClientOptions()
	opt.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.MqttHost, cfg.MqttPort))
	opt.SetClientID("ttrn")
	opt.SetUsername(cfg.MqttUser)
	opt.SetPassword(cfg.MqttPassword)
	opt.OnConnect = rsl.connectHandler
	opt.OnConnectionLost = rsl.connectionLostHandler
	opt.SetDefaultPublishHandler(rsl.defaultHandler)

	rsl.broker = mqtt.NewClient(opt)
	rsl.broker.Subscribe("/la", 0, func(client mqtt.Client, msg mqtt.Message) {
		logrus.Warnf("got new unhandled mqtt message on topic %s, with payload %s", msg.Topic(), msg.Payload())
	})
	rsl.broker.Subscribe("/train", 1, rsl.trainHandler)
	return rsl
}

// Run connects to the MQTT broker.
func (m Mqtt) Run() {
	if token := m.broker.Connect(); token.Wait() && token.Error() != nil {
		logrus.Panicf("failed to connect to broker: %s", token.Error())
	}
}

func (m Mqtt) defaultHandler(client mqtt.Client, msg mqtt.Message) {
	logrus.Warnf("got new unhandled mqtt message on topic %s, with payload %s", msg.Topic(), msg.Payload())
}

func (m Mqtt) connectHandler(client mqtt.Client) {
	logrus.Info("mqtt connected")
}

func (m Mqtt) connectionLostHandler(client mqtt.Client, err error) {
	logrus.Errorf("mqtt connection lost %s", err)
}

func (m Mqtt) trainHandler(client mqtt.Client, msg mqtt.Message) {
	logrus.Debugf("got train update on topic %s with payload %s", msg.Topic(), msg.Payload())
	logrus.Info("hoi")
	// speedRe := regexp.MustCompile(`/train/(\d{1})/speed`)
}
