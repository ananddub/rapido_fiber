package repo

import (
	"sync"
	"testing"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func TestMQTTPublishSubscribe(t *testing.T) {

	var wg sync.WaitGroup
	wg.Add(1)

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("test-client-1")
	opts.SetCleanSession(true)

	messageReceived := make(chan string)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		messageReceived <- string(msg.Payload())
		wg.Done()
	})

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Connect error: %v", token.Error())
	}

	// Subscribe
	if token := client.Subscribe("test/topic", 1, nil); token.Wait() && token.Error() != nil {
		t.Fatalf("Subscribe error: %v", token.Error())
	}
	// Publish
	token := client.Publish("test/topic", 1, false, "hello-test")
	token.Wait()

	// Wait for message
	select {
	case msg := <-messageReceived:
		if msg != "hello-test" {
			t.Fatalf("Expected hello-test, got %s", msg)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Timeout waiting for message")
	}

	wg.Wait()
	client.Disconnect(250)
}
