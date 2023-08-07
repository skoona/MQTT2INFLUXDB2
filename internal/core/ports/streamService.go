package ports

type StreamService interface {
	Enable() error
	Disable()
	GetDeviceRepo() StorageRepository
	GetStreamConsumer() StreamConsumer
	GetStreamProvider() StreamProvider
	ChartEnvironmentals(msg StreamMessage)
}
