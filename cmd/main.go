package main

import (
	"context"
	"fmt"
	"mqttToInfluxDB/internal/services"
	"os"
	"os/signal"
	"syscall"
)

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
	map[string]Nodes
			   map[string]Properties


*/

func main() {
	ctxService, cancelService := context.WithCancel(context.Background())

	service := services.NewStreamService(ctxService)

	service.Enable()

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
