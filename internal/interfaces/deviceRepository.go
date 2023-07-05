package interfaces

import (
	"fyne.io/fyne/v2/data/binding"
	"mqttToInfluxDB/internal/entities"
)

type DeviceRepository interface {
	ApplyMessage(msg StreamMessage)
	NewDevice(msg StreamMessage) *entities.Device

	GetNamedDevice(deviceName string) *entities.Device
	GetNamedProperty(deviceName, property string) *entities.Property
	GetDevices() map[string]*entities.Device
	GetProperties(deviceName string) map[string]*entities.Property
	GetMessageCount() *binding.ExternalString
	GetDeviceCount() *binding.ExternalString
}
