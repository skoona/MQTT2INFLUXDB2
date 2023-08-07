package ports

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/skoona/mqttToInfluxDB/internal/core/domain"
)

type StorageRepository interface {
	ApplyMessage(msg StreamMessage)
	NewDevice(msg StreamMessage) *domain.Device

	GetNamedDevice(deviceName string) *domain.Device
	GetNamedProperty(deviceName, property string) *domain.Property
	GetDevices() map[string]*domain.Device
	GetProperties(deviceName string) map[string]*domain.Property
	GetMessageCount() *binding.ExternalString
	GetDeviceCount() *binding.ExternalString
}
