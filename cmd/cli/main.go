//go:build cli
// +build cli

package main

/*
 * MONITORS
 * - sknSensors/+/Ambient/humidity float
 * - sknSensors/+/Ambient/temperature float
 * - sknSensors/+/Occupancy/motion  ON,OFF
 * - sknSensors/+/Occupancy/occupancy ON,OFF
 * - sknSensors/+/Presence/motion OPEN,CLOSED
 * - sknSensors/+/SknRanger/Position int
 * - sknSensors/+/SknRanger/State UP,DOWN
 * - sknSensors/+/SknRanger/Details JSON

	Device / Node / Property / Value
	map[string]Device
			   map[string]Properties

*/

import (
	"context"
	"fmt"
	"mqttToInfluxDB/internal/commons"
	"mqttToInfluxDB/internal/services"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, commons.SknAppIDKey, "net.skoona.mq2influx")

	_ = commons.AppSettings(ctx)

	ctx = context.WithValue(ctx, commons.InfluxHostUriKey, commons.GetInfluxHostUri()) // strings
	ctx = context.WithValue(ctx, commons.InfluxBucketKey, commons.GetInfluxBucket())
	ctx = context.WithValue(ctx, commons.InfluxOrgKey, commons.GetInfluxOrg())
	ctx = context.WithValue(ctx, commons.InfluxTokenKey, commons.GetInfluxToken())
	ctx = context.WithValue(ctx, commons.MqttHostUriKey, commons.GetMqttHostUri())
	ctx = context.WithValue(ctx, commons.MqttUserKey, commons.GetMqttUser())
	ctx = context.WithValue(ctx, commons.MqttPassKey, commons.GetMqttPass())
	ctx = context.WithValue(ctx, commons.DebugModeKey, commons.IsDebugMode()) // bool
	ctx = context.WithValue(ctx, commons.TestModeKey, commons.IsTestMode())   // bool
	ctxService, cancelService := context.WithCancel(ctx)

	enbledDataStore := false
	service := services.NewStreamService(ctxService, commons.IsInfluxDBEnabled(), enbledDataStore, nil)
	if err := service.Enable(); err != nil {
		fmt.Println("ERROR: shutdown requested cause:", err.Error())
		fmt.Println("Likely a configuration error, use environment vars to set runtime config values")
		cancelService()

	} else {
		/*
		 * Prepare for clean exit
		 */
		errs := make(chan error, 1)
		go func(shutdown chan error) {
			systemSignalChannel := make(chan os.Signal, 1)
			signal.Notify(systemSignalChannel, syscall.SIGINT, syscall.SIGTERM)
			sig := <-systemSignalChannel // wait on ctrl-c
			cancelService()              // provider

			shutdown <- fmt.Errorf("%s", sig)
		}(errs)
		fmt.Println("event ", "shutdown requested ", "cause:", <-errs) // errs holds it up
	}

	time.Sleep(3 * time.Second)
}
