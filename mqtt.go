package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const topicName = "tool/accesscontrol"

var c MQTT.Client

func connectToMQTT() {
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
	opts.SetClientID("tool-access-control")

	//create and start the client 
	c = MQTT.NewClient(opts)
	c.Connect()
}

func publish(message string) {
	c.Publish(topicName, 0, false, message)
}
