package ui

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"mqttToInfluxDB/internal/entities"
	"mqttToInfluxDB/internal/interfaces"
)

type ViewProvider interface {
	UpdateUI() bool
	MainPage() *fyne.Container
	NewCard(device entities.Device) *widget.Card
}

type viewProvider struct {
	ctx      context.Context
	cards    map[string]*widget.Card
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
		cards:   map[string]*widget.Card{},
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

func (v *viewProvider) NewCard(device entities.Device) *widget.Card {
	device.SetDisplayed(true)
	props := widget.NewForm()
	for name, prop := range device.Properties {
		item := widget.NewFormItem(name, widget.NewLabel(prop.Value))
		props.AppendItem(item)
	}
	card := widget.NewCard(device.Name, "example", props)
	v.cards[device.Name] = card

	return card
}
func (v *viewProvider) UpdateUI() bool {
	for _, dev := range v.service.GetDeviceRepo().GetDevices() {
		v.NewCard(dev)
	}
	v.MainPage()
	return true
}
func (v *viewProvider) MainPage() *fyne.Container {
	msg := fmt.Sprintf("Device Count: %v", len(v.service.GetDeviceRepo().GetDevices()))
	v.status.SetText(msg)
	grid := container.NewGridWrap(fyne.NewSize(300, 200))
	for _, card := range v.cards {
		grid.Add(card)
	}
	v.mainPage = container.NewBorder(nil, container.NewHBox(v.refresh, v.status), nil, nil, grid)
	return v.mainPage
}
