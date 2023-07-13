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
	chart           sknlinechart.SknLineChart
}

var _ interfaces.StreamService = (*streamService)(nil)

func NewStreamService(ctx context.Context, enableInflux bool, enabledDataStore bool, linechart sknlinechart.SknLineChart) interfaces.StreamService {
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

	point := sknlinechart.NewLineChartDatapoint(float32(val), string(clr), time.RFC3339)
	s.chart.ApplyDataPoint(series, &point)
}
