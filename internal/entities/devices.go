package entities

type (
	Device struct {
		Name       string
		DeviceType string
		Displayed  bool
		Properties map[string]Property
	}
	Property struct {
		Name  string
		Value string
	}
)

func (d *Device) IsDisplayed() bool {
	return d.Displayed
}
func (d *Device) SetDisplayed(onOff bool) {
	d.Displayed = onOff
}
