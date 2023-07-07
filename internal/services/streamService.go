package services

import (
	"context"
	"fmt"
	"mqttToInfluxDB/internal/commons"
	"mqttToInfluxDB/internal/interfaces"
	"mqttToInfluxDB/internal/repositories"
)

type streamService struct {
	ctx             context.Context
	enableDataStore bool
	enableInflux    bool
	stream          chan interfaces.StreamMessage
	provider        interfaces.StreamProvider
	consumer        interfaces.StreamConsumer
	devStore        interfaces.StreamStorage
}

var _ interfaces.StreamService = (*streamService)(nil)

func NewStreamService(ctx context.Context, enableInflux bool, enabledDataStore bool) interfaces.StreamService {
	var consumer interfaces.StreamConsumer
	var devStore interfaces.StreamStorage

	stream := make(chan interfaces.StreamMessage, 64)
	if enabledDataStore {
		devStore = repositories.NewStreamStorage(ctx)
	}
	if enableInflux {
		consumer = repositories.NewStreamConsumer(ctx)
	}
	provider := repositories.NewStreamProvider(ctx, stream)

	return &streamService{
		ctx:             ctx,
		enableDataStore: enabledDataStore,
		enableInflux:    enableInflux,
		stream:          stream,
		provider:        provider,
		consumer:        consumer,
		devStore:        devStore,
	}
}

func (s *streamService) Enable() error {

	go func(ctx context.Context, stream chan interfaces.StreamMessage, consume interfaces.StreamConsumer, devStore interfaces.StreamStorage) {
		debug := ctx.Value(commons.DebugModeKey).(bool)
		fmt.Println("====> StreamService() Listening")
		for msg := range stream {
			if debug {
				fmt.Printf("[%s] DEVICE: %s\tPROPERTY: %s VALUE: %v\n", msg.Timestamp(), msg.Device(), msg.Property(), msg.Value())
			}
			if s.enableDataStore {
				devStore.ApplyMessage(msg)
			}
			if s.enableInflux {
				if msg.Property() != "heartbeat" {
					_ = consume.Write(msg)
				}
			}
		}
	}(s.ctx, s.stream, s.consumer, s.devStore)

	err := s.provider.Connect()
	if err != nil {
		fmt.Println("Connect() error", err.Error())
		return err
	}
	err = s.provider.EnableStream()
	if err != nil {
		fmt.Println("EnableStream() error", err.Error())
		return err
	}
	return nil
}
func (s *streamService) Disable() {
	s.provider.DisableStream()
	close(s.stream)
	//cancelConsumer() // consumer
	//cancelDevice()   // devStore

}
func (s *streamService) GetStreamProvider() interfaces.StreamProvider {
	return s.provider
}
func (s *streamService) GetStreamConsumer() interfaces.StreamConsumer {
	return s.consumer
}
func (s *streamService) GetDeviceRepo() interfaces.StreamStorage {
	return s.devStore
}
