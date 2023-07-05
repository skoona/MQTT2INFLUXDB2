package repositories

import (
	"context"
	"crypto/tls"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"mqttToInfluxDB/internal/commons"
	"mqttToInfluxDB/internal/interfaces"
	"time"
)

type consumer struct {
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
	ctx      context.Context
	stream   chan interfaces.StreamMessage
}

var _ interfaces.StreamConsumer = (*consumer)(nil)

func NewStreamConsumer(ctx context.Context, stream chan interfaces.StreamMessage, devStore interfaces.DeviceRepository) interfaces.StreamConsumer {

	debug := ctx.Value(commons.DebugModeKey).(bool)

	bucket := ctx.Value(commons.InfluxBucketKey).(string)
	org := ctx.Value(commons.InfluxOrgKey).(string)
	token := ctx.Value(commons.InfluxTokenKey).(string)
	url := ctx.Value(commons.InfluxHostUriKey).(string)

	repo := &consumer{ctx: ctx, stream: stream}

	repo.client = influxdb2.NewClientWithOptions(url, token,
		influxdb2.DefaultOptions().
			SetUseGZip(true).
			SetTLSConfig(&tls.Config{
				InsecureSkipVerify: true,
			}))

	repo.writeAPI = repo.client.WriteAPIBlocking(org, bucket)

	go func(consume interfaces.StreamConsumer, devStore interfaces.DeviceRepository) {
		fmt.Println("====> StreamConsumer() Listening")
		for msg := range consume.GetStream() {
			if debug {
				fmt.Printf("[%s] DEVICE: %s\tPROPERTY: %s VALUE: %v\n", msg.Timestamp(), msg.Device(), msg.Property(), msg.Value())
			}
			devStore.ApplyMessage(msg)
			_ = consume.Write(msg)
		}
	}(repo, devStore)

	go func(ctx context.Context, consume interfaces.StreamConsumer) {
		for {
			if <-ctx.Done(); true {
				fmt.Println("consumer cancelled\n", ctx.Err())
				consume.Disconnect()
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}(ctx, repo)

	return repo
}
func (c *consumer) GetStream() chan interfaces.StreamMessage {
	return c.stream
}
func (c *consumer) Write(msg interfaces.StreamMessage) error {
	if msg.IsGarageDoor() {
		dataPoint := influxdb2.NewPoint("home",
			map[string]string{
				"device": msg.Device(),
				"node":   msg.Node(),
			},
			map[string]interface{}{
				"position": msg.Position(),
				"actual":   msg.Actual(),
				"state":    msg.State(),
				"ambient":  msg.Ambient(),
				"signal":   msg.SignalStrength(),
			}, time.Now(),
		)
		err := c.writeAPI.WritePoint(c.ctx, dataPoint)
		if err != nil {
			fmt.Println("error while sending to influx: ", err.Error())
			return err
		}
	} else {
		dataPoint := influxdb2.NewPoint("home",
			map[string]string{
				"device": msg.Device(),
				"node":   msg.Node(),
			},
			map[string]interface{}{msg.Property(): msg.Value()}, time.Now(),
		)
		err := c.writeAPI.WritePoint(c.ctx, dataPoint)
		if err != nil {
			fmt.Println("error while sending to influx: ", err.Error())
			return err
		}
	}
	return nil
}
func (c *consumer) Disconnect() {
	fmt.Println("====> StreamConsumer() disconnected")
	_ = c.writeAPI.Flush(c.ctx)
	c.client.Close()
}
