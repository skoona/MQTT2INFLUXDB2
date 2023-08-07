package ports

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/skoona/mqttToInfluxDB/internal/core/domain"
)

type StreamService interface {
	Enable() error
	Disable()
	IsStreamProviderEnabled() bool
	IsStreamConsumerEnabled() bool
	GetDeviceList() map[string]*domain.Device
	GetMessageCount() *binding.String
	GetDeviceCount() *binding.String
	ChartEnvironmentals(msg StreamMessage)
}
