package ports

type StreamConsumer interface {
	ApplyMessage(msg StreamMessage) error
	Disconnect()
}
