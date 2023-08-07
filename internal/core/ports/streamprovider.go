package ports

type StreamProvider interface {
	IsOnline() bool
	Connect() error
	Disconnect()
	EnableStream() error
	DisableStream() error
}
