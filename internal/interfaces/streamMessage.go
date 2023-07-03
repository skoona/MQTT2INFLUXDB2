package interfaces

type StreamMessage interface {
	IsGarageDoor() bool
	Topic() string
	Network() string
	Device() string
	Node() string
	Property() string
	Value() string
	Timestamp() string
	Actual() int
	Ambient() float32
	Position() int
	SignalStrength() float32
	State() string
}
