package repositories

import (
	"context"
	"fmt"
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
	dType := entities.SensorType
	if msg.IsGarageDoor() {
		dType = entities.GarageType
	}
	device := entities.Device{
		Name:       msg.Device(),
		DeviceType: dType,
		Properties: map[string]entities.Property{},
	}
	device.Properties[msg.Property()] = entities.Property{
		Name:  msg.Property(),
		Value: msg.Value(),
	}

	if msg.IsGarageDoor() {
		device.Properties[entities.ActualProperty] = entities.Property{
			Name:  entities.ActualProperty,
			Value: fmt.Sprint(msg.Actual()),
		}
		device.Properties[entities.AmbientProperty] = entities.Property{
			Name:  entities.AmbientProperty,
			Value: fmt.Sprint(msg.Ambient()),
		}
		device.Properties[entities.PositionProperty] = entities.Property{
			Name:  entities.PositionProperty,
			Value: fmt.Sprint(msg.Position()),
		}
		device.Properties[entities.SignalStrengthProperty] = entities.Property{
			Name:  entities.SignalStrengthProperty,
			Value: fmt.Sprint(msg.SignalStrength()),
		}
		device.Properties[entities.StateProperty] = entities.Property{
			Name:  entities.StateProperty,
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
		property, ok := device.Properties[entities.ActualProperty]
		if !ok {
			device.Properties[entities.ActualProperty] = entities.Property{
				Name:  entities.ActualProperty,
				Value: fmt.Sprint(msg.Actual()),
			}
		}
		property.Value = fmt.Sprint(msg.Actual())

		property, ok = device.Properties[entities.AmbientProperty]
		if !ok {
			device.Properties[entities.AmbientProperty] = entities.Property{
				Name:  entities.AmbientProperty,
				Value: fmt.Sprint(msg.Ambient()),
			}
		}
		property.Value = fmt.Sprint(msg.Ambient())

		property, ok = device.Properties[entities.PositionProperty]
		if !ok {
			device.Properties[entities.PositionProperty] = entities.Property{
				Name:  entities.PositionProperty,
				Value: fmt.Sprint(msg.Position()),
			}
		}
		property.Value = fmt.Sprint(msg.Position())

		property, ok = device.Properties[entities.SignalStrengthProperty]
		if !ok {
			device.Properties[entities.SignalStrengthProperty] = entities.Property{
				Name:  entities.SignalStrengthProperty,
				Value: fmt.Sprint(msg.SignalStrength()),
			}
		}
		property.Value = fmt.Sprint(msg.SignalStrength())

		property, ok = device.Properties[entities.StateProperty]
		if !ok {
			device.Properties[entities.StateProperty] = entities.Property{
				Name:  entities.StateProperty,
				Value: fmt.Sprint(msg.State()),
			}
		}
		property.Value = fmt.Sprint(msg.State())

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
