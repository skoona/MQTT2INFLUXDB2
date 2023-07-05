package interfaces

type StreamService interface {
	Enable() error
	Disable()
	GetDeviceRepo() DeviceRepository
	GetStreamConsumer() StreamConsumer
	GetStreamProvider() StreamProvider
}
