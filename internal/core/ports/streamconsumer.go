package ports

type StreamConsumer interface {
	Write(msg StreamMessage) error
	Disconnect()
}
