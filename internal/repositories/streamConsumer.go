package repositories

import (
	"context"
	"crypto/tls"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"mqttToInfluxDB/internal/interfaces"
	"os"
	"time"
)

type consumer struct {
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
	ctx      context.Context
	stream   chan interfaces.StreamMessage
}

var _ interfaces.StreamConsumer = (*consumer)(nil)

func NewStreamConsumer(ctx context.Context, stream chan interfaces.StreamMessage) interfaces.StreamConsumer {

	bucket := os.Getenv("INFLUXDB_BUCKET")
	org := os.Getenv("INFLUXDB_ORG")
	token := os.Getenv("INFLUXDB_TOKEN")
	url := os.Getenv("INFLUXDB_URI")

	repo := &consumer{ctx: ctx, stream: stream}

	repo.client = influxdb2.NewClientWithOptions(url, token,
		influxdb2.DefaultOptions().
			SetUseGZip(true).
			SetTLSConfig(&tls.Config{
				InsecureSkipVerify: true,
			}))

	repo.writeAPI = repo.client.WriteAPIBlocking(org, bucket)

	go func(consume interfaces.StreamConsumer) {
		fmt.Println("====> StreamConsumer() Listening")
		for msg := range consume.GetStream() {
			//fmt.Printf("[%s] DEVICE: %s NODE: %s \tRECEIVED PROPERTY: %s VALUE: %v\n", msg.Timestamp(), msg.Device(), msg.Node(), msg.Property(), msg.Value())
			_ = consume.Write(msg)
		}
	}(repo)

	go func(ctx context.Context, consume interfaces.StreamConsumer) {
		for true {
			select {
			case <-ctx.Done():
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
	c.client.Close()
}
