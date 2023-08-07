package repository

import (
	"context"
	"crypto/tls"
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/skoona/mqttToInfluxDB/internal/commons"
	"github.com/skoona/mqttToInfluxDB/internal/core/ports"
	"time"
)

type influxConsumer struct {
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
	ctx      context.Context
}

var _ ports.StreamConsumer = (*influxConsumer)(nil)

func NewStreamConsumer(ctx context.Context) ports.StreamConsumer {

	bucket := ctx.Value(commons.InfluxBucketKey).(string)
	org := ctx.Value(commons.InfluxOrgKey).(string)
	token := ctx.Value(commons.InfluxTokenKey).(string)
	url := ctx.Value(commons.InfluxHostUriKey).(string)

	repo := &influxConsumer{ctx: ctx}

	repo.client = influxdb2.NewClientWithOptions(url, token,
		influxdb2.DefaultOptions().
			SetUseGZip(true).
			SetTLSConfig(&tls.Config{
				InsecureSkipVerify: true,
			}))

	repo.writeAPI = repo.client.WriteAPIBlocking(org, bucket)

	go func(ctx context.Context, consume ports.StreamConsumer) {
		for {
			if <-ctx.Done(); true {
				fmt.Println("influxConsumer cancelled\n", ctx.Err())
				consume.Disconnect()
				break
			}
			time.Sleep(1 * time.Second)
		}
	}(ctx, repo)

	return repo
}
func (c *influxConsumer) ApplyMessage(msg ports.StreamMessage) error {
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
func (c *influxConsumer) Disconnect() {
	fmt.Println("====> StreamConsumer() disconnected")
	_ = c.writeAPI.Flush(c.ctx)
	c.client.Close()
}
