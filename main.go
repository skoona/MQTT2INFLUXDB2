//go:build gui
// +build gui

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
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/skoona/mqttToInfluxDB/internal/commons"
	"github.com/skoona/mqttToInfluxDB/internal/services"
	"github.com/skoona/mqttToInfluxDB/internal/ui"
	"github.com/skoona/sknlinechart"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, commons.SknAppIDKey, "net.skoona.mq2influx")

	gui := app.NewWithID("net.skoona.mq2influx")
	win := gui.NewWindow("MQTT to InfluxDB2")
	lgw := gui.NewWindow("Line Graph")
	_ = commons.AppSettings(ctx)

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

	points := map[string][]*sknlinechart.ChartDatapoint{}
	skn, err := sknlinechart.NewLineChart("Skoona's Home Automation Network", "Inside mqtt 2 influx", &points)
	if err != nil {
		fmt.Println(err.Error())
	}
	skn.SetMiddleLeftLabel("Temperature")
	skn.SetMiddleRightLabel("Humidity")
	skn.SetBottomLeftLabel("sknSensors MQTT Network")

	onLine := true
	enbledDataStore := true
	service := services.NewStreamService(ctxService, commons.IsInfluxDBEnabled(), enbledDataStore, skn)
	err = service.Enable()
	if err != nil {
		// configuration failure
		onLine = false
	}

	sknMenus(gui, win)
	SknTrayMenu(gui, win, lgw)
	win.Resize(fyne.NewSize(1024, 756))

	viewProvider := ui.NewViewProvider(ctxService, service)
	if onLine {
		time.Sleep(3 * time.Second)
		win.SetContent(viewProvider.MainPage())

		lgw.Resize(fyne.NewSize(982, 452))
		lgw.SetContent(container.NewPadded(skn))
		lgw.CenterOnScreen()
		lgw.SetCloseIntercept(func() { lgw.Hide() })

		win.Show()
		lgw.Show()
	} else {
		if err != nil {
			win.SetContent(viewProvider.ConfigFailedPage(err.Error()))
		} else {
			win.SetContent(viewProvider.ConfigFailedPage("Unknown Error"))
		}
		win.Show()
	}

	go func() {
		systemSignalChannel := make(chan os.Signal, 1)
		signal.Notify(systemSignalChannel, syscall.SIGINT, syscall.SIGTERM)
		sig := <-systemSignalChannel // wait on ctrl-c
		cancelService()              // provider
		fmt.Println(sig.String())
		gui.Quit()
	}()

	gui.Run()

	cancelService() // provider
	time.Sleep(3 * time.Second)

}
