package services

import (
	"context"
	"fmt"
	"mqttToInfluxDB/internal/commons"
	"mqttToInfluxDB/internal/interfaces"
	"mqttToInfluxDB/internal/repositories"
)

type streamService struct {
	ctx      context.Context
	stream   chan interfaces.StreamMessage
	provider interfaces.StreamProvider
	consumer interfaces.StreamConsumer
	devStore interfaces.DeviceRepository
}

var _ interfaces.StreamService = (*streamService)(nil)

func NewStreamService(ctx context.Context, enableInflux bool) interfaces.StreamService {
	var consumer interfaces.StreamConsumer

	stream := make(chan interfaces.StreamMessage, 64)
	devStore := repositories.NewDeviceRepository(ctx)
	provider := repositories.NewStreamProvider(ctx, stream)
	if enableInflux {
		consumer = repositories.NewStreamConsumer(ctx)
	}

	return &streamService{
		ctx:      ctx,
		stream:   stream,
		provider: provider,
		consumer: consumer,
		devStore: devStore,
	}
}

func (s *streamService) Enable() error {

	go func(ctx context.Context, stream chan interfaces.StreamMessage, consume interfaces.StreamConsumer, devStore interfaces.DeviceRepository) {
		debug := ctx.Value(commons.DebugModeKey).(bool)
		fmt.Println("====> StreamService() Listening")
		for msg := range stream {
			if debug {
				fmt.Printf("[%s] DEVICE: %s\tPROPERTY: %s VALUE: %v\n", msg.Timestamp(), msg.Device(), msg.Property(), msg.Value())
			}
			devStore.ApplyMessage(msg)
			if consume != nil {
				_ = consume.Write(msg)
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
func (s *streamService) GetDeviceRepo() interfaces.DeviceRepository {
	return s.devStore
}
