package ttrn

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
)

// Mqtt bundles the MQTT related stuff.
type Mqtt struct {
	broker        mqtt.Client
	turnoutPositionEvents chan TurnoutPositionEvent
	trainSpeedEvents   chan TrainSpeedEvent
	trainPositionEvents chan TrainPositionEvent
}

// NewMqtt returns a new Mqtt instance.
func NewMqtt(cfg Config, turnoutPositionEvents chan TurnoutPositionEvent, trainSpeedEvents chan TrainSpeedEvent, trainPositionEvents chan TrainPositionEvent) Mqtt {
	rsl := Mqtt{
		turnoutPositionEvents: turnoutPositionEvents,
		trainSpeedEvents:   trainSpeedEvents,
		trainPositionEvents: trainPositionEvents,
	}

	opt := mqtt.NewClientOptions()
	opt.AddBroker(fmt.Sprintf("tcp://%s:%d", cfg.MqttHost, cfg.MqttPort))
	opt.SetClientID("ttrn")
	opt.SetUsername(cfg.MqttUser)
	opt.SetPassword(cfg.MqttPassword)
	opt.OnConnect = rsl.connectHandler
	opt.OnConnectionLost = rsl.connectionLostHandler
	// opt.SetDefaultPublishHandler(rsl.defaultHandler)

	rsl.broker = mqtt.NewClient(opt)
	return rsl
}

// Run connects to the MQTT broker.
func (m Mqtt) Run() {
	if token := m.broker.Connect(); token.Wait() && token.Error() != nil {
		logrus.Errorf("failed to connect to broker: %s", token.Error())
		os.Exit(1)
	}
	m.init()
	m.broker.Subscribe("/train/+/position", 1, m.trainPositionHandler)
}

// init defaults all topics.
func (m Mqtt) init() {
	m.publish("/train/0/speed", "0")
	m.publish("/train/1/speed", "0")
	m.publish("/train/2/speed", "0")
	m.publish("/turnout/0/position", "0")
	m.publish("/turnout/1/position", "0")
	m.publish("/turnout/2/position", "0")
	m.publish("/turnout/3/position", "0")
	m.publish("/turnout/4/position", "0")
	m.publish("/turnout/5/position", "0")
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

func (m Mqtt) publish(topic string, payload interface{}) {
	if token := m.broker.Publish(topic, 1, true, payload); token.Wait() && token.Error() != nil {
		logrus.Errorf("couldn't publish to topic %s with payload %+v", topic, payload)
	}
}

func (m Mqtt) trainPositionHandler(client mqtt.Client, msg mqtt.Message) {
	logrus.Debugf("got train position update on topic %s with payload %s", msg.Topic(), msg.Payload())
	positionRe := regexp.MustCompile(`/train/(\d{1})/position`)
	rsl := positionRe.FindAllStringSubmatch(msg.Topic(), -1)
	if len(rsl) != 1 && len(rsl[0]) !=2 {
		logrus.Errorf("received topic %s has a illegal format for a train position topic", msg.Topic())
		return
	}
	id, err := strconv.Atoi(rsl[0][1])
	if err != nil {
		logrus.Errorf("error while converting the given train id %s to int", rsl[1])
		return
	}
	if id < 0 || id > 2 {
		logrus.Errorf("given train id %d isn't valid", id)
		return
	}

	position, err := strconv.Atoi(string(msg.Payload()))
	if err != nil {
		logrus.Errorf("error while converting position value %s to int", msg.Payload())
		return
	}
	if position < 0 || position > 3 {
		logrus.Errorf("received train-position %d isn't valid", position)
	}
	
	logrus.Debugf("received new position %d for train %d", position, id)
}
