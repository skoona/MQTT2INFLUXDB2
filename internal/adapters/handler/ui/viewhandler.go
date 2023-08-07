package ui

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skoona/mqttToInfluxDB/internal/commons"
	"github.com/skoona/mqttToInfluxDB/internal/core/domain"
	"github.com/skoona/mqttToInfluxDB/internal/core/ports"
	"sync"
)

type ViewHandler interface {
	UpdateUI() bool
	MainPage() *fyne.Container
	ConfigFailedPage(msg string) *fyne.Container
	NewCard(device *domain.Device) *fyne.Container
	SetStatusLineText(msg string)
}

type viewHandler struct {
	ctx          context.Context
	cards        map[string]*fyne.Container
	mainPageGrid *fyne.Container
	mainPage     *fyne.Container
	status       *widget.Label
	refresh      *widget.Button
	msgCounter   *widget.Label
	devCounter   *widget.Label
	mainWindow   fyne.Window
	service      ports.StreamService
	updateUiLock sync.RWMutex
}

var _ ViewHandler = (*viewHandler)(nil)

func NewViewHandler(ctx context.Context, service ports.StreamService) ViewHandler {
	win := ctx.Value(commons.FyneWindowKey).(*fyne.Window)
	view := &viewHandler{
		ctx:          ctx,
		service:      service,
		mainWindow:   *win,
		cards:        map[string]*fyne.Container{},
		status:       widget.NewLabel("place holder"),
		mainPageGrid: container.NewGridWithColumns(4),
		updateUiLock: sync.RWMutex{},
	}

	view.msgCounter = widget.NewLabelWithData(*view.service.GetMessageCount())
	view.devCounter = widget.NewLabelWithData(*view.service.GetDeviceCount())
	(*service.GetMessageCount()).AddListener(binding.NewDataListener(func() {
		if len(view.cards) > 1 {
			if view.UpdateUI() {
				view.mainWindow.Content().Refresh()
			}
		}
	}))

	view.refresh = widget.NewButtonWithIcon("refresh", theme.ViewRefreshIcon(), func() {
		if view.UpdateUI() {
			view.mainWindow.Content().Refresh()
		}
	})

	if view.UpdateUI() {
		view.mainWindow.Content().Refresh()
	}

	return view
}
func (v *viewHandler) NewCard(device *domain.Device) *fyne.Container {
	device.SetDisplayed(true)
	border := canvas.NewRectangle(theme.BackgroundColor())
	border.StrokeColor = theme.InputBorderColor()
	border.StrokeWidth = 4
	props := container.New(layout.NewFormLayout())
	for name, prop := range device.Properties {
		if name != commons.GarageProperty {
			z := widget.NewLabel(prop.Name)
			z.Alignment = fyne.TextAlignTrailing
			props.Add(z)
			n := widget.NewLabelWithData(prop.Bond)
			n.Wrapping = fyne.TextWrapWord
			props.Add(n)
		}
	}
	card := widget.NewCard(device.Name, device.UpdatedAt(), props)
	device.Bond = binding.BindString(&device.LastUpdate)
	device.Bond.AddListener(binding.NewDataListener(func() {
		str, _ := device.Bond.Get()
		card.SetSubTitle(str)
	}))

	if device.IsGarageType() {
		if device.IsGarageOpen() {
			card.SetImage(commons.SknSelectThemedImage("garageOpen"))
		} else {
			card.SetImage(commons.SknSelectThemedImage("garageClosed"))
		}
	} else {
		card.SetImage(commons.SknSelectThemedImage("sensorOn_o"))
	}
	card.Resize(fyne.NewSize(100, 100))
	content := container.NewMax(border, card)
	v.cards[device.Name] = content
	v.mainPageGrid.Add(content)
	v.mainPageGrid.Refresh()

	return content
}
func (v *viewHandler) UpdateUI() bool {
	v.updateUiLock.Lock()
	defer v.updateUiLock.Unlock()

	added := false
	for _, dev := range v.service.GetDeviceList() {
		if !dev.IsDisplayed() {
			v.NewCard(dev)
			v.SetStatusLineText("added new Device: " + dev.Name)
			added = true
		} else {
			// only update properties not on screen
			card, ok := v.cards[dev.Name]
			if ok {
				if dev.IsGarageType() {
					if dev.IsGarageOpen() {
						card.Objects[1].(*widget.Card).SetImage(commons.SknSelectThemedImage("garageOpen"))
					} else {
						card.Objects[1].(*widget.Card).SetImage(commons.SknSelectThemedImage("garageClosed"))
					}

				} else {
					for name, prop := range dev.Properties {
						if name != commons.GarageProperty {
							skip := false
							// find in card then skip
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
								d.Wrapping = fyne.TextWrapWord

								card.Objects[1].(*widget.Card).Content.(*fyne.Container).Add(n)
								card.Objects[1].(*widget.Card).Content.(*fyne.Container).Add(d)
								added = true
								v.SetStatusLineText("added new Property: " + dev.Name + "::" + prop.Name)
							}
						}
					}
				}
			} // device didn't exist
		}
	}

	return added
}
func (v *viewHandler) MainPage() *fyne.Container {
	v.SetStatusLineText("main page updated")
	v.mainPageGrid.RemoveAll()
	for _, card := range v.cards {
		v.mainPageGrid.Add(card)
	}
	m := widget.NewIcon(theme.FolderOpenIcon())
	i := widget.NewIcon(theme.StorageIcon())
	if !v.service.IsStreamProviderEnabled() {
		m.Hide()
	}
	if !v.service.IsStreamConsumerEnabled() {
		i.Hide()
	}

	scrolledGrid := container.NewVScroll(v.mainPageGrid)

	v.mainPage = container.NewBorder(
		nil,
		container.NewHBox(v.refresh,
			m,
			i,
			widget.NewLabel(" Devices:"), v.devCounter,
			widget.NewLabel(" Messages processed:"), v.msgCounter,
			v.status,
		),
		nil,
		nil,
		scrolledGrid)

	return v.mainPage
}
func (v *viewHandler) SetStatusLineText(msg string) {
	v.status.SetText(msg)
}
func (v *viewHandler) ConfigFailedPage(msg string) *fyne.Container {
	title := canvas.NewText("Configuration Failure", theme.PrimaryColor())
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Italic: true}
	title.TextSize = 24

	eLine := canvas.NewText(msg, theme.ErrorColor())
	eLine.Alignment = fyne.TextAlignCenter
	eLine.TextSize = 18

	body := canvas.NewText("set run configuration in menu `settings`", theme.WarningColor())
	body.Alignment = fyne.TextAlignCenter
	body.TextSize = 18

	dialog.ShowError(fmt.Errorf("configurating error: %s", fmt.Errorf(msg)), v.mainWindow)

	return container.NewMax(
		container.NewVBox(title),
		container.NewCenter(
			container.NewVBox(layout.NewSpacer(), eLine, body),
		),
	)
}
