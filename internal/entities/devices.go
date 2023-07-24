package entities

import (
	"fyne.io/fyne/v2/data/binding"
	"github.com/skoona/mqttToInfluxDB/internal/commons"
)

type (
	Device struct {
		Name       string
		DeviceType string
		Displayed  bool
		LastUpdate string
		Bond       binding.ExternalString // to LastUpdate
		Properties map[string]*Property
	}
	Property struct {
		Name  string
		Bond  binding.ExternalString // to Value
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
	if state, ok := d.Properties[commons.StateProperty]; ok && (state.Value == "UP" || state.Value == "OPEN") {
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
