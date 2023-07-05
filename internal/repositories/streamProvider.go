package repositories

import (
	"context"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/go-uuid"
	"mqttToInfluxDB/internal/commons"
	"mqttToInfluxDB/internal/entities"
	"mqttToInfluxDB/internal/interfaces"
	"time"
)

var (
	topics = []string{}

	topicsMap = map[string]byte{
		"sknSensors/+/Ambient/humidity":    byte(1),
		"sknSensors/+/Ambient/temperature": byte(1),
		"sknSensors/+/Occupancy/motion":    byte(1),
		"sknSensors/+/Occupancy/occupancy": byte(1),
		"sknSensors/+/Presence/motion":     byte(1),
		"sknSensors/+/SknRanger/Position":  byte(1),
		"sknSensors/+/SknRanger/State":     byte(1),
		"sknSensors/+/SknRanger/Details":   byte(1),
	}
)

type repo struct {
	client     MQTT.Client
	subscribed bool
	stream     chan interfaces.StreamMessage
}

var (
	_            interfaces.StreamProvider = (*repo)(nil)
	mqttProvider *repo
)

func NewStreamProvider(ctx context.Context, stream chan interfaces.StreamMessage) interfaces.StreamProvider {
	for k := range topicsMap {
		topics = append(topics, k)
	}

	clientIdValue, _ := uuid.GenerateUUID()
	opts := MQTT.NewClientOptions()
	opts.AddBroker(ctx.Value(commons.MqttHostUriKey).(string))
	opts.SetClientID(clientIdValue)
	opts.SetUsername(ctx.Value(commons.MqttUserKey).(string))
	opts.SetPassword(ctx.Value(commons.MqttPassKey).(string))
	opts.SetCleanSession(true)
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		sm, _ := entities.NewStreamMessage(msg.Topic(), string(msg.Payload()))
		stream <- sm
	})
	opts.OnConnect = func(client MQTT.Client) {
		fmt.Println("====> StreamProvider() Connected")
	}
	opts.OnReconnecting = func(client MQTT.Client, options *MQTT.ClientOptions) {
		fmt.Println("====> StreamProvider() Reconnecting")
	}
	opts.OnConnectionLost = func(client MQTT.Client, err error) {
		fmt.Println("====> StreamProvider() Connection lost: ", err.Error())
	}

	client := MQTT.NewClient(opts)

	mqttProvider = &repo{
		client:     client,
		subscribed: false,
		stream:     stream,
	}

	go func(ctx context.Context, provider interfaces.StreamProvider) {
		for {
			if <-ctx.Done(); true {
				fmt.Println("provider cancelled\n", ctx.Err())
				provider.DisableStream()
				provider.Disconnect()
				break
			}
			time.Sleep(1 * time.Second)
		}
	}(ctx, mqttProvider)

	return mqttProvider
}
func GetClient() MQTT.Client {
	return mqttProvider.client
}
func (r *repo) IsOnline() bool {
	return r.client.IsConnected()
}
func (r *repo) Connect() error {
	fmt.Println("====> StreamProvider() Connecting")
	if token := r.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
func (r *repo) Disconnect() {
	fmt.Println("====> StreamProvider() Disconnecting")
	r.client.Disconnect(250)
}
func (r *repo) EnableStream() error {
	for !r.IsOnline() {
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("====> StreamProvider() Subscribing")
	token := r.client.SubscribeMultiple(topicsMap, func(client MQTT.Client, msg MQTT.Message) {
		sm, _ := entities.NewStreamMessage(msg.Topic(), string(msg.Payload()))
		r.stream <- sm
	})
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return token.Error()
	}

	r.subscribed = true
	return nil
}
func (r *repo) DisableStream() error {
	if !r.subscribed {
		return nil
	}
	fmt.Println("====> StreamProvider() Unsubscribing")
	token := r.client.Unsubscribe(topics...)
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return token.Error()
	}
	r.subscribed = false
	return nil
}
