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
	"mqttToInfluxDB/internal/services"
	"mqttToInfluxDB/internal/ui"
)

func main() {
	ctxService, cancelService := context.WithCancel(context.Background())
	service := services.NewStreamService(ctxService)
	service.Enable()

	gui := app.New()
	win := gui.NewWindow("MQTT to InfluxDB2")
	win.Resize(fyne.NewSize(1024, 756))

	viewProvider := ui.NewViewProvider(ctxService, service)
	win.SetContent(viewProvider.MainPage())

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
