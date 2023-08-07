package services

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2/theme"
	"github.com/skoona/mqttToInfluxDB/internal/adapters/repository"
	"github.com/skoona/mqttToInfluxDB/internal/commons"
	"github.com/skoona/mqttToInfluxDB/internal/core/ports"
	"github.com/skoona/sknlinechart"
	"strconv"
	"time"
)

type streamService struct {
	ctx             context.Context
	enableDataStore bool
	enableInflux    bool
	stream          chan ports.StreamMessage
	provider        ports.StreamProvider
	consumer        ports.StreamConsumer
	devStore        ports.StorageRepository
	chart           sknlinechart.LineChart
}

var _ ports.StreamService = (*streamService)(nil)

func NewStreamService(ctx context.Context, enableInflux bool, enabledDataStore bool, linechart sknlinechart.LineChart) ports.StreamService {
	var consumer ports.StreamConsumer
	var devStore ports.StorageRepository

	stream := make(chan ports.StreamMessage, 64)
	if enabledDataStore {
		devStore = repository.NewStorageRepository(ctx)
	}
	if enableInflux {
		consumer = repository.NewStreamConsumer(ctx)
	}
	provider := repository.NewStreamProvider(ctx, stream)

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
			if s.chart != nil {
				s.ChartEnvironmentals(msg)
			}
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
	_ = s.provider.DisableStream()
	close(s.stream)
}
func (s *streamService) GetStreamProvider() ports.StreamProvider {
	return s.provider
}
func (s *streamService) GetStreamConsumer() ports.StreamConsumer {
	return s.consumer
}
func (s *streamService) GetDeviceRepo() ports.StorageRepository {
	return s.devStore
}

func (s *streamService) ChartEnvironmentals(msg ports.StreamMessage) {
	if msg.Property() != "temperature" && msg.Property() != "humidity" && msg.Property() != "Position" {
		return
	}

	val, _ := strconv.ParseFloat(msg.Value(), 32)
	clr := theme.ColorNameForeground

	// [red orange yellow green blue purple brown gray]
	series := msg.Device() + "::" + msg.Property()

	switch msg.Device() {
	case "GarageMonitor":
		if msg.Property() == "temperature" {
			clr = theme.ColorRed
		} else {
			clr = theme.ColorBrown
		}
	case "GuestRoom":
		if msg.Property() == "temperature" {
			clr = theme.ColorOrange
		} else {
			clr = theme.ColorGray
		}
	case "FamilyRoom":
		if msg.Property() == "temperature" {
			clr = theme.ColorYellow
		} else {
			clr = theme.ColorRed
		}
	case "OutsideMonitor":
		if msg.Property() == "temperature" {
			clr = theme.ColorGreen
		} else {
			clr = theme.ColorYellow
		}
	case "MediaRoom":
		if msg.Property() == "temperature" {
			clr = theme.ColorBlue
		} else {
			clr = theme.ColorGreen
		}
	case "HomeOffice":
		if msg.Property() == "temperature" {
			clr = theme.ColorPurple
		} else {
			clr = theme.ColorBlue
		}
	case "OverheadDoor":
		if msg.Property() == "Position" {
			clr = theme.ColorNameForeground
		}
	default:
		clr = theme.ColorNameForeground
	}

	point := sknlinechart.NewChartDatapoint(float32(val), string(clr), time.RFC3339)
	s.chart.ApplyDataPoint(series, &point)
}
