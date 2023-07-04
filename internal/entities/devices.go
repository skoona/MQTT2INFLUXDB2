package entities

type (
	Device struct {
		Name       string
		DeviceType string
		Properties map[string]Property
	}
	Property struct {
		Name  string
		Value string
	}
)
