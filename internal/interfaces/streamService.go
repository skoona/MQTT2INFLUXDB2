package interfaces

type StreamService interface {
	Enable()
	Disable()
	GetDeviceRepo() DeviceRepository
	GetStreamConsumer() StreamConsumer
	GetStreamProvider() StreamProvider
}
