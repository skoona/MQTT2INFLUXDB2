package ports

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/skoona/mqttToInfluxDB/internal/core/domain"
)

type StorageRepository interface {
	ApplyMessage(msg StreamMessage)
	NewDevice(msg StreamMessage) *domain.Device
	GetDevices() map[string]*domain.Device
	GetMessageCountBinding() *binding.String
	GetDeviceCountBinding() *binding.String
}
