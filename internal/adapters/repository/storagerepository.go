package repository

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2/data/binding"
	"github.com/skoona/mqttToInfluxDB/internal/commons"
	"github.com/skoona/mqttToInfluxDB/internal/core/domain"
	"github.com/skoona/mqttToInfluxDB/internal/core/ports"
	"strconv"
)

type storageRepository struct {
	devices    map[string]*domain.Device
	ctx        context.Context
	msgCount   int
	bMsgCntStr string
	bMsgCount  binding.ExternalString
	bDevCntStr string
	bDevCount  binding.ExternalString
}

func NewStorageRepository(ctx context.Context) ports.StorageRepository {
	devices := &storageRepository{
		devices: map[string]*domain.Device{},
		ctx:     ctx,
	}
	devices.bMsgCount = binding.BindString(&devices.bMsgCntStr)
	devices.bDevCount = binding.BindString(&devices.bDevCntStr)

	return devices
}

func (d *storageRepository) NewDevice(msg ports.StreamMessage) *domain.Device {
	dType := commons.SensorType
	if msg.IsGarageDoor() {
		dType = commons.GarageType
	}

	device := &domain.Device{
		Name:       msg.Device(),
		DeviceType: dType,
		LastUpdate: msg.Timestamp(),
		Properties: map[string]*domain.Property{},
	}
	device.Properties[msg.Property()] = &domain.Property{
		Name:  msg.Property(),
		Value: msg.Value(),
	}
	prop := device.Properties[msg.Property()]
	prop.Bond = binding.BindString(&prop.Value)

	if msg.IsGarageDoor() {
		device.Properties[commons.ActualProperty] = &domain.Property{
			Name:  commons.ActualProperty,
			Value: fmt.Sprintf("%d", msg.Actual()),
		}
		prop = device.Properties[commons.ActualProperty]
		prop.Bond = binding.BindString(&prop.Value)

		device.Properties[commons.AmbientProperty] = &domain.Property{
			Name:  commons.AmbientProperty,
			Value: fmt.Sprintf("%3.2f", msg.Ambient()),
		}
		prop = device.Properties[commons.AmbientProperty]
		prop.Bond = binding.BindString(&prop.Value)

		device.Properties[commons.PositionProperty] = &domain.Property{
			Name:  commons.PositionProperty,
			Value: fmt.Sprintf("%d", msg.Position()),
		}
		prop = device.Properties[commons.PositionProperty]
		prop.Bond = binding.BindString(&prop.Value)

		device.Properties[commons.SignalStrengthProperty] = &domain.Property{
			Name:  commons.SignalStrengthProperty,
			Value: fmt.Sprintf("%3.2f", msg.SignalStrength()),
		}
		prop = device.Properties[commons.SignalStrengthProperty]
		prop.Bond = binding.BindString(&prop.Value)

		device.Properties[commons.StateProperty] = &domain.Property{
			Name:  commons.StateProperty,
			Value: msg.State(),
		}
		prop = device.Properties[commons.StateProperty]
		prop.Bond = binding.BindString(&prop.Value)

	}

	d.devices[msg.Device()] = device

	d.bDevCntStr = strconv.Itoa(len(d.devices))
	_ = d.bDevCount.Set(d.bDevCntStr)

	return device
}
func (d *storageRepository) ApplyMessage(msg ports.StreamMessage) {
	device, ok := d.devices[msg.Device()]
	d.msgCount += 1
	d.bMsgCntStr = strconv.Itoa(d.msgCount)
	_ = d.bMsgCount.Set(d.bMsgCntStr)

	if !ok {
		_ = d.NewDevice(msg)
		return
	}
	device.LastUpdate = msg.Timestamp()
	if device.Bond != nil {
		_ = device.Bond.Reload()
	}

	var prop *domain.Property

	property, ok := device.Properties[msg.Property()]
	if !ok {
		device.Properties[msg.Property()] = &domain.Property{
			Name:  msg.Property(),
			Value: msg.Value(),
		}
		prop = device.Properties[msg.Property()]
		prop.Bond = binding.BindString(&prop.Value)

	} else {
		property.Value = msg.Value()
		_ = property.Bond.Set(property.Value)
	}

	if msg.IsGarageDoor() {
		device.DeviceType = commons.GarageType

		property, ok := device.Properties[commons.ActualProperty]
		if !ok {
			device.Properties[commons.ActualProperty] = &domain.Property{
				Name:  commons.ActualProperty,
				Value: fmt.Sprintf("%d", msg.Actual()),
			}
			prop = device.Properties[commons.ActualProperty]
			prop.Bond = binding.BindString(&prop.Value)
		} else {
			property.Value = fmt.Sprintf("%d", msg.Actual())
			_ = property.Bond.Set(property.Value)
		}

		property, ok = device.Properties[commons.AmbientProperty]
		if !ok {
			device.Properties[commons.AmbientProperty] = &domain.Property{
				Name:  commons.AmbientProperty,
				Value: fmt.Sprintf("%3.2f", msg.Ambient()),
			}
			prop = device.Properties[commons.AmbientProperty]
			prop.Bond = binding.BindString(&prop.Value)
		} else {
			property.Value = fmt.Sprintf("%3.2f", msg.Ambient())
			_ = property.Bond.Set(property.Value)
		}

		property, ok = device.Properties[commons.PositionProperty]
		if !ok {
			device.Properties[commons.PositionProperty] = &domain.Property{
				Name:  commons.PositionProperty,
				Value: fmt.Sprintf("%d", msg.Position()),
			}
			prop = device.Properties[commons.PositionProperty]
			prop.Bond = binding.BindString(&prop.Value)
		} else {
			property.Value = fmt.Sprintf("%d", msg.Position())
			_ = property.Bond.Set(property.Value)
		}

		property, ok = device.Properties[commons.SignalStrengthProperty]
		if !ok {
			device.Properties[commons.SignalStrengthProperty] = &domain.Property{
				Name:  commons.SignalStrengthProperty,
				Value: fmt.Sprintf("%3.2f", msg.SignalStrength()),
			}
			prop = device.Properties[commons.SignalStrengthProperty]
			prop.Bond = binding.BindString(&prop.Value)
		} else {
			property.Value = fmt.Sprintf("%3.2f", msg.SignalStrength())
			_ = property.Bond.Set(property.Value)
		}

		property, ok = device.Properties[commons.StateProperty]
		if !ok {
			device.Properties[commons.StateProperty] = &domain.Property{
				Name:  commons.StateProperty,
				Value: msg.State(),
			}
			prop = device.Properties[commons.StateProperty]
			prop.Bond = binding.BindString(&prop.Value)
		} else {
			property.Value = msg.State()
			_ = property.Bond.Set(property.Value)
		}
	}

}
func (d *storageRepository) GetNamedDevice(deviceName string) *domain.Device {
	return d.devices[deviceName]
}
func (d *storageRepository) GetNamedProperty(deviceName, property string) *domain.Property {
	return d.devices[deviceName].Properties[property]
}
func (d *storageRepository) GetDevices() map[string]*domain.Device {
	return d.devices
}
func (d *storageRepository) GetProperties(deviceName string) map[string]*domain.Property {
	return d.devices[deviceName].Properties
}
func (d *storageRepository) GetMessageCount() *binding.ExternalString {
	return &d.bMsgCount
}
func (d *storageRepository) GetDeviceCount() *binding.ExternalString {
	return &d.bDevCount
}
