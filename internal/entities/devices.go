package entities

import "mqttToInfluxDB/internal/commons"

type (
	Device struct {
		Name       string
		DeviceType string
		Displayed  bool
		LastUpdate string
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
func (d *Device) IsGarageType() bool {
	return d.DeviceType == commons.GarageType
}
func (d *Device) IsGarageOpen() bool {
	if state, ok := d.Properties[commons.PositionProperty]; ok && (state.Value == "UP" || state.Value == "OPEN") {
		return true
	}
	return false
}
func (d *Device) SetDisplayed(onOff bool) {
	d.Displayed = onOff
}
func (d *Device) UpdatedAt() string {
	return d.LastUpdate
}
