package interfaces

type StreamService interface {
	Enable() error
	Disable()
	GetDeviceRepo() StreamStorage
	GetStreamConsumer() StreamConsumer
	GetStreamProvider() StreamProvider
}
