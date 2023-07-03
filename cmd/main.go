package main

import (
	"context"
	"fmt"
	"mqttToInfluxDB/internal/interfaces"
	"mqttToInfluxDB/internal/repositories"
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
 */

func main() {
	stream := make(chan interfaces.StreamMessage, 64)

	ctxProvider, cancelProvider := context.WithCancel(context.Background())
	ctxConsumer, cancelConsumer := context.WithCancel(context.Background())

	provider, err := repositories.NewStreamProvider(ctxProvider, stream)
	if err != nil {
		fmt.Println("NewStreamProvider() error", err.Error())
		panic(1)
	}

	err = provider.Connect()
	if err != nil {
		fmt.Println("Connect() error", err.Error())
		panic(1)
	}

	_ = repositories.NewStreamConsumer(ctxConsumer, stream)

	err = provider.EnableStream()
	if err != nil {
		fmt.Println("EnableStream() error", err.Error())
		panic(1)
	}

	/*
	 * Prepare for clean exit
	 */
	errs := make(chan error, 1)
	go func(shutdown chan error) {
		systemSignalChannel := make(chan os.Signal, 1)
		signal.Notify(systemSignalChannel, syscall.SIGINT, syscall.SIGTERM)
		sig := <-systemSignalChannel // wait on ctrl-c

		cancelProvider() // provider
		close(stream)
		cancelConsumer() // consumer

		shutdown <- fmt.Errorf("%s", sig)
	}(errs)
	fmt.Println("event ", "shutdown requested ", "cause:", <-errs)

}
