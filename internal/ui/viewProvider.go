package ui

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
	ConfigFailedPage(msg string) *fyne.Container
	NewCard(device *entities.Device) *fyne.Container
	SetStatusLineText(msg string)
}

type viewProvider struct {
	ctx        context.Context
	cards      map[string]*fyne.Container
	mainPage   *fyne.Container
	status     *widget.Label
	refresh    *widget.Button
	mainWindow fyne.Window
	service    interfaces.StreamService
}

var _ ViewProvider = (*viewProvider)(nil)

func NewViewProvider(ctx context.Context, service interfaces.StreamService) ViewProvider {
	win := ctx.Value(commons.FyneWindowKey).(*fyne.Window)
	view := &viewProvider{
		ctx:        ctx,
		service:    service,
		mainWindow: *win,
		cards:      map[string]*fyne.Container{},
		status:     widget.NewLabel("place holder"),
	}
	view.refresh = widget.NewButtonWithIcon("refresh", theme.ViewRefreshIcon(), func() {
		msg := fmt.Sprintf("Device Count: %v", len(view.service.GetDeviceRepo().GetDevices()))
		view.status.SetText(msg)
		if view.UpdateUI() {
			view.mainWindow.SetContent(view.mainPage)
		} else {
			view.mainWindow.Content().Refresh()
		}
	})

	view.UpdateUI()
	return view
}
func (v *viewProvider) NewCard(device *entities.Device) *fyne.Container {
	device.SetDisplayed(true)
	border := canvas.NewRectangle(theme.OverlayBackgroundColor()) //  .InputBackgroundColor())
	border.StrokeColor = theme.InputBorderColor()
	border.StrokeWidth = 4
	props := container.New(layout.NewFormLayout())
	for name, prop := range device.Properties {
		if name != commons.GarageProperty {
			//if prop.Bond == nil {
			//	prop.Bond = binding.BindString(&prop.Value)
			//}
			n := widget.NewLabel(prop.Name)
			d := widget.NewLabelWithData(prop.Bond)

			props.Add(n)
			props.Add(d)
		}
	}
	card := widget.NewCard(device.Name, device.UpdatedAt(), props)
	device.Bond = binding.BindString(&device.LastUpdate)
	callback := binding.NewDataListener(func() {
		str, _ := device.Bond.Get()
		card.SetSubTitle(str)
	})
	device.Bond.AddListener(callback)

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
	added := false
	for _, dev := range v.service.GetDeviceRepo().GetDevices() {
		if !dev.IsDisplayed() {
			v.NewCard(dev)
			v.mainWindow.SetContent(v.MainPage())
			v.SetStatusLineText("added new Device: " + dev.Name)
			added = true
		} else {
			// only update properties not on screen
			card, _ := v.cards[dev.Name]
			if dev.IsGarageType() {
				if dev.IsGarageOpen() {
					card.Objects[1].(*widget.Card).SetImage(commons.SknSelectThemedImage("garageOpen"))
				} else {
					card.Objects[1].(*widget.Card).SetImage(commons.SknSelectThemedImage("garageClosed"))
				}
			}

			for name, prop := range dev.Properties {
				if name != commons.GarageProperty {
					skip := false
					// find in card then add
					for _, wid := range card.Objects[1].(*widget.Card).Content.(*fyne.Container).Objects {
						item := wid.(*widget.Label)
						if item.Text == name {
							skip = true
							break
						}
					}
					if !skip {
						n := widget.NewLabel(prop.Name)
						d := widget.NewLabelWithData(prop.Bond)
						card.Objects[1].(*widget.Card).Content.(*fyne.Container).Add(n)
						card.Objects[1].(*widget.Card).Content.(*fyne.Container).Add(d)
					}
				}
			}
		}
	}
	return added
}
func (v *viewProvider) MainPage() *fyne.Container {
	msg := fmt.Sprintf("Device Count: %v", len(v.service.GetDeviceRepo().GetDevices()))
	v.SetStatusLineText(msg)
	grid := container.NewGridWithColumns(4)
	for _, card := range v.cards {
		grid.Add(card)
	}
	v.mainPage = container.NewBorder(nil, container.NewHBox(v.refresh, v.status), nil, nil, grid)

	return v.mainPage
}
func (v *viewProvider) SetStatusLineText(msg string) {
	v.status.SetText(msg)
}
func (v *viewProvider) ConfigFailedPage(msg string) *fyne.Container {
	title := canvas.NewText("Configuration Failure", theme.PrimaryColor())
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Italic: true}
	title.TextSize = 24

	eLine := canvas.NewText(msg, theme.ErrorColor())
	eLine.Alignment = fyne.TextAlignCenter
	eLine.TextSize = 18

	body := canvas.NewText("Set run configuration in menu `settings`", theme.WarningColor())
	body.Alignment = fyne.TextAlignCenter
	body.TextSize = 18

	return container.NewMax(
		container.NewVBox(title),
		container.NewCenter(
			container.NewVBox(layout.NewSpacer(), eLine, body),
		),
	)
}
