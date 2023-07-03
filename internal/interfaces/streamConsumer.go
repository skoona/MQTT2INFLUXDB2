package interfaces

type StreamConsumer interface {
	Write(msg StreamMessage) error
	Disconnect()
	GetStream() chan StreamMessage
}
