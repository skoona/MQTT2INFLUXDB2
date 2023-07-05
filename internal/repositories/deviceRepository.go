package repositories

import (
	"context"
	"fmt"
	"mqttToInfluxDB/internal/commons"
	"mqttToInfluxDB/internal/entities"
	"mqttToInfluxDB/internal/interfaces"
)

type deviceProvider struct {
	devices map[string]entities.Device
	ctx     context.Context
}

func NewDeviceRepository(ctx context.Context) interfaces.DeviceRepository {
	devices := &deviceProvider{
		devices: map[string]entities.Device{},
		ctx:     ctx,
	}

	return devices
}

func (d *deviceProvider) NewDevice(msg interfaces.StreamMessage) entities.Device {
	dType := commons.SensorType
	if msg.IsGarageDoor() {
		dType = commons.GarageType
	}
	device := entities.Device{
		Name:       msg.Device(),
		DeviceType: dType,
		LastUpdate: msg.Timestamp(),
		Properties: map[string]entities.Property{},
	}
	device.Properties[msg.Property()] = entities.Property{
		Name:  msg.Property(),
		Value: msg.Value(),
	}

	if msg.IsGarageDoor() {
		device.Properties[commons.ActualProperty] = entities.Property{
			Name:  commons.ActualProperty,
			Value: fmt.Sprintf("%d", msg.Actual()),
		}
		device.Properties[commons.AmbientProperty] = entities.Property{
			Name:  commons.AmbientProperty,
			Value: fmt.Sprintf("%3.2f", msg.Ambient()),
		}
		device.Properties[commons.PositionProperty] = entities.Property{
			Name:  commons.PositionProperty,
			Value: fmt.Sprintf("%d", msg.Position()),
		}
		device.Properties[commons.SignalStrengthProperty] = entities.Property{
			Name:  commons.SignalStrengthProperty,
			Value: fmt.Sprintf("%3.2f", msg.SignalStrength()),
		}
		device.Properties[commons.StateProperty] = entities.Property{
			Name:  commons.StateProperty,
			Value: msg.State(),
		}
	}

	d.devices[msg.Device()] = device

	return device
}
func (d *deviceProvider) ApplyMessage(msg interfaces.StreamMessage) {
	device, ok := d.devices[msg.Device()]
	if !ok {
		device = d.NewDevice(msg)
		return
	}
	device.LastUpdate = msg.Timestamp()

	property, ok := device.Properties[msg.Property()]
	if !ok {
		device.Properties[msg.Property()] = entities.Property{
			Name:  msg.Property(),
			Value: msg.Value(),
		}
	} else {
		property.Value = msg.Value()
	}

	if msg.IsGarageDoor() {
		device.DeviceType = commons.GarageType

		property, ok := device.Properties[commons.ActualProperty]
		if !ok {
			device.Properties[commons.ActualProperty] = entities.Property{
				Name:  commons.ActualProperty,
				Value: fmt.Sprintf("%d", msg.Actual()),
			}
		}
		property.Value = fmt.Sprintf("%d", msg.Actual())

		property, ok = device.Properties[commons.AmbientProperty]
		if !ok {
			device.Properties[commons.AmbientProperty] = entities.Property{
				Name:  commons.AmbientProperty,
				Value: fmt.Sprintf("%3.2f", msg.Ambient()),
			}
		}
		property.Value = fmt.Sprintf("%3.2f", msg.Ambient())

		property, ok = device.Properties[commons.PositionProperty]
		if !ok {
			device.Properties[commons.PositionProperty] = entities.Property{
				Name:  commons.PositionProperty,
				Value: fmt.Sprintf("%d", msg.Position()),
			}
		}
		property.Value = fmt.Sprintf("%d", msg.Position())

		property, ok = device.Properties[commons.SignalStrengthProperty]
		if !ok {
			device.Properties[commons.SignalStrengthProperty] = entities.Property{
				Name:  commons.SignalStrengthProperty,
				Value: fmt.Sprintf("%3.2f", msg.SignalStrength()),
			}
		}
		property.Value = fmt.Sprintf("%3.2f", msg.SignalStrength())

		property, ok = device.Properties[commons.StateProperty]
		if !ok {
			device.Properties[commons.StateProperty] = entities.Property{
				Name:  commons.StateProperty,
				Value: msg.State(),
			}
		}
		property.Value = msg.State()
	}

}
func (d *deviceProvider) GetNamedDevice(deviceName string) entities.Device {
	return d.devices[deviceName]
}
func (d *deviceProvider) GetNamedProperty(deviceName, property string) entities.Property {
	return d.devices[deviceName].Properties[property]
}
func (d *deviceProvider) GetDevices() map[string]entities.Device {
	return d.devices
}
func (d *deviceProvider) GetProperties(deviceName string) map[string]entities.Property {
	return d.devices[deviceName].Properties
}
