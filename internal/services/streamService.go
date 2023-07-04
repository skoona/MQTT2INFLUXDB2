package services

import (
	"context"
	"fmt"
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

func NewStreamService(ctx context.Context) interfaces.StreamService {
	stream := make(chan interfaces.StreamMessage, 64)
	devStore := repositories.NewDeviceRepository(ctx)
	provider := repositories.NewStreamProvider(ctx, stream)
	consumer := repositories.NewStreamConsumer(ctx, stream, devStore)

	return &streamService{
		ctx:      ctx,
		stream:   stream,
		provider: provider,
		consumer: consumer,
		devStore: devStore,
	}
}

func (s *streamService) Enable() {
	err := s.provider.Connect()
	if err != nil {
		fmt.Println("Connect() error", err.Error())
		panic(1)
	}
	err = s.provider.EnableStream()
	if err != nil {
		fmt.Println("EnableStream() error", err.Error())
		panic(1)
	}

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
