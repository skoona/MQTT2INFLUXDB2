package repository

import (
	"context"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/hashicorp/go-uuid"
	"github.com/skoona/mqttToInfluxDB/internal/commons"
	"github.com/skoona/mqttToInfluxDB/internal/core/domain"
	"github.com/skoona/mqttToInfluxDB/internal/core/ports"
	"time"
)

var (
	topics = []string{}

	topicsMap = map[string]byte{
		"+/+/+/humidity":    byte(1),
		"+/+/+/temperature": byte(1),
		"+/+/+/motion":      byte(1),
		"+/+/+/occupancy":   byte(1),
		"+/+/+/Position":    byte(1),
		"+/+/+/State":       byte(1),
		"+/+/+/Details":     byte(1),
		"+/+/+/message":     byte(1),
		"+/+/+/name":        byte(1),
		"+/+/+/heartbeat":   byte(1),
	}
)

type homieProvider struct {
	client     MQTT.Client
	subscribed bool
	stream     chan ports.StreamMessage
}

var (
	_            ports.StreamProvider = (*homieProvider)(nil)
	mqttProvider *homieProvider
)

func NewStreamProvider(ctx context.Context, stream chan ports.StreamMessage) ports.StreamProvider {
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
		sm, _ := domain.NewStreamMessage(msg.Topic(), string(msg.Payload()))
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

	mqttProvider = &homieProvider{
		client:     client,
		subscribed: false,
		stream:     stream,
	}

	go func(ctx context.Context, provider ports.StreamProvider) {
		for {
			if <-ctx.Done(); true {
				fmt.Println("provider cancelled\n", ctx.Err())
				_ = provider.DisableStream()
				provider.Disconnect()
				break
			}
			time.Sleep(1 * time.Second)
		}
	}(ctx, mqttProvider)

	return mqttProvider
}
func (r *homieProvider) IsOnline() bool {
	return r.client.IsConnected()
}
func (r *homieProvider) Connect() error {
	fmt.Println("====> StreamProvider() Connecting")
	if token := r.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}
func (r *homieProvider) Disconnect() {
	fmt.Println("====> StreamProvider() Disconnecting")
	r.client.Disconnect(250)
}
func (r *homieProvider) EnableStream() error {
	for !r.IsOnline() {
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("====> StreamProvider() Subscribing")
	token := r.client.SubscribeMultiple(topicsMap, func(client MQTT.Client, msg MQTT.Message) {
		sm, _ := domain.NewStreamMessage(msg.Topic(), string(msg.Payload()))
		r.stream <- sm
	})
	if token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		return token.Error()
	}

	r.subscribed = true
	return nil
}
func (r *homieProvider) DisableStream() error {
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
