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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"mqttToInfluxDB/internal/commons"
	"mqttToInfluxDB/internal/services"
	"mqttToInfluxDB/internal/ui"
	"time"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, commons.SknAppIDKey, "net.skoona.mq2influx")

	gui := app.NewWithID("net.skoona.mq2influx")
	win := gui.NewWindow("MQTT to InfluxDB2")
	_ = commons.AppSettings(ctx, gui)

	ctx = context.WithValue(ctx, commons.InfluxHostUriKey, commons.GetInfluxHostUri()) // strings
	ctx = context.WithValue(ctx, commons.InfluxBucketKey, commons.GetInfluxBucket())
	ctx = context.WithValue(ctx, commons.InfluxOrgKey, commons.GetInfluxOrg())
	ctx = context.WithValue(ctx, commons.InfluxTokenKey, commons.GetInfluxToken())
	ctx = context.WithValue(ctx, commons.MqttHostUriKey, commons.GetMqttHostUri())
	ctx = context.WithValue(ctx, commons.MqttUserKey, commons.GetMqttUser())
	ctx = context.WithValue(ctx, commons.MqttPassKey, commons.GetMqttPass())
	ctx = context.WithValue(ctx, commons.FyneWindowKey, &win)
	ctx = context.WithValue(ctx, commons.DebugModeKey, commons.IsDebugMode()) // bool
	ctx = context.WithValue(ctx, commons.TestModeKey, commons.IsTestMode())   // bool
	ctxService, cancelService := context.WithCancel(ctx)

	onLine := true
	service := services.NewStreamService(ctxService)
	err := service.Enable()
	if err != nil {
		// configuration failure
		err.Error()
		onLine = false
	}

	sknMenus(gui, win)
	SknTrayMenu(gui, win)
	win.Resize(fyne.NewSize(1024, 756))

	time.Sleep(4 * time.Second)

	viewProvider := ui.NewViewProvider(ctxService, service)
	if onLine {
		win.SetContent(viewProvider.MainPage())
	} else {
		if err != nil {
			win.SetContent(viewProvider.ConfigFailedPage(err.Error()))
		} else {
			win.SetContent(viewProvider.ConfigFailedPage("Unknown Error"))
		}
	}
	win.ShowAndRun()
	cancelService() // provider

	/*
		 * Prepare for clean exit
		errs := make(chan error, 1)
		go func(shutdown chan error) {
			systemSignalChannel := make(chan os.Signal, 1)
			signal.Notify(systemSignalChannel, syscall.SIGINT, syscall.SIGTERM)
			sig := <-systemSignalChannel // wait on ctrl-c
			cancelService()              // provider

			shutdown <- fmt.Errorf("%s", sig)
		}(errs)
		fmt.Println("event ", "shutdown requested ", "cause:", <-errs) // errs holds it up
	*/

}
