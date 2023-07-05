package interfaces

type StreamConsumer interface {
	Write(msg StreamMessage) error
	Disconnect()
}
