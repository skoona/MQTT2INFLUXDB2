package ui

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"mqttToInfluxDB/internal/commons"
	"mqttToInfluxDB/internal/entities"
	"mqttToInfluxDB/internal/interfaces"
)

type ViewProvider interface {
	UpdateUI() bool
	MainPage() *fyne.Container
	NewCard(device entities.Device) *fyne.Container
}

type viewProvider struct {
	ctx      context.Context
	cards    map[string]*fyne.Container
	mainPage *fyne.Container
	status   *widget.Label
	refresh  *widget.Button
	service  interfaces.StreamService
}

var _ ViewProvider = (*viewProvider)(nil)

func NewViewProvider(ctx context.Context, service interfaces.StreamService) ViewProvider {
	view := &viewProvider{
		ctx:     ctx,
		service: service,
		cards:   map[string]*fyne.Container{},
		status:  widget.NewLabel("place holder"),
	}
	view.refresh = widget.NewButtonWithIcon("refresh", theme.ViewRefreshIcon(), func() {
		msg := fmt.Sprintf("Device Count: %v", len(view.service.GetDeviceRepo().GetDevices()))
		view.status.SetText(msg)
		view.UpdateUI()
	})

	view.UpdateUI()
	return view
}

func (v *viewProvider) NewCard(device entities.Device) *fyne.Container {
	device.SetDisplayed(true)
	border := canvas.NewRectangle(theme.InputBackgroundColor())
	border.StrokeColor = theme.InputBorderColor()
	border.StrokeWidth = 6
	props := container.New(layout.NewFormLayout())
	for name, prop := range device.Properties {
		if name != commons.GarageProperty {
			n := widget.NewLabel(prop.Name)
			v := widget.NewLabel(prop.Value)

			props.Add(n)
			props.Add(v)
		}
	}
	card := widget.NewCard(device.Name, device.UpdatedAt(), props)
	if device.IsGarageType() {
		if device.IsGarageOpen() {
			card.SetImage(commons.SknSelectThemedImage("garageOpen"))
		} else {
			card.SetImage(commons.SknSelectThemedImage("garageClosed"))
		}
	} else {
		card.SetImage(commons.SknSelectThemedImage("sensorOn_o"))
	}
	content := container.NewMax(border, card)
	content.Resize(fyne.NewSize(240, 288))
	v.cards[device.Name] = content

	return content
}
func (v *viewProvider) UpdateUI() bool {
	for _, dev := range v.service.GetDeviceRepo().GetDevices() {
		v.NewCard(dev)
	}
	v.MainPage()
	v.mainPage.Refresh()
	return true
}
func (v *viewProvider) MainPage() *fyne.Container {
	msg := fmt.Sprintf("Device Count: %v", len(v.service.GetDeviceRepo().GetDevices()))
	v.status.SetText(msg)
	grid := container.NewGridWithColumns(4)
	for _, card := range v.cards {
		grid.Add(card)
	}
	v.mainPage = container.NewBorder(nil, container.NewHBox(v.refresh, v.status), nil, nil, grid)
	return v.mainPage
}
