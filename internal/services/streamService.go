package services

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2/theme"
	"github.com/skoona/sknlinechart"
	"mqttToInfluxDB/internal/commons"
	"mqttToInfluxDB/internal/interfaces"
	"mqttToInfluxDB/internal/repositories"
	"strconv"
	"time"
)

type streamService struct {
	ctx             context.Context
	enableDataStore bool
	enableInflux    bool
	stream          chan interfaces.StreamMessage
	provider        interfaces.StreamProvider
	consumer        interfaces.StreamConsumer
	devStore        interfaces.StreamStorage
	chart           sknlinechart.LineChart
}

var _ interfaces.StreamService = (*streamService)(nil)

func NewStreamService(ctx context.Context, enableInflux bool, enabledDataStore bool, linechart sknlinechart.LineChart) interfaces.StreamService {
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
		chart:           linechart,
	}
}

func (s *streamService) Enable() error {

	go func(svc *streamService) {
		debug := svc.ctx.Value(commons.DebugModeKey).(bool)
		fmt.Println("====> StreamService() Listening")
		for msg := range svc.stream {
			if debug {
				fmt.Printf("[%s] DEVICE: %s\tPROPERTY: %s VALUE: %v\n", msg.Timestamp(), msg.Device(), msg.Property(), msg.Value())
			}
			if s.enableDataStore {
				svc.devStore.ApplyMessage(msg)
			}
			if s.enableInflux {
				if msg.Property() != "heartbeat" {
					_ = svc.consumer.Write(msg)
				}
			}
			s.ChartEnvironmentals(msg)
		}
	}(s)

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

func (s *streamService) ChartEnvironmentals(msg interfaces.StreamMessage) {
	if msg.Property() != "temperature" && msg.Property() != "humidity" && msg.Property() != "Position" {
		return
	}

	series := msg.Device() + "::" + msg.Property()

	val, _ := strconv.ParseFloat(msg.Value(), 32)
	clr := theme.ColorNameForeground
	if msg.Property() == "Position" {
		clr = theme.ColorPurple
	}
	if msg.Property() == "temperature" {
		clr = theme.ColorYellow
	}
	if msg.Property() == "humidity" {
		clr = theme.ColorBlue
	}

	point := sknlinechart.NewLineChartDatapoint(float32(val), string(clr), time.RFC3339)
	s.chart.ApplyDataPoint(series, point)
}
