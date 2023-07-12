//go:build gui
// +build gui

package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"mqttToInfluxDB/internal/commons"
	"net/url"
)

func PreferencesPage() *fyne.Container {
	desc := canvas.NewText("Runtime Configuration", color.White)
	desc.Alignment = fyne.TextAlignCenter
	desc.TextStyle = fyne.TextStyle{Italic: true}
	desc.TextSize = 24

	influxUrl := widget.NewEntry()
	influxUrl.SetPlaceHolder("InfluxDB2 ip http://10.100.1.17:8086")
	influxUrl.SetText(commons.GetInfluxHostUri())
	influxBucket := widget.NewEntry()
	influxBucket.SetPlaceHolder("InfluxDB2 bucket SknSensors")
	influxBucket.SetText(commons.GetInfluxBucket())
	influxOrg := widget.NewEntry()
	influxOrg.SetPlaceHolder("InfluxDB2 Organization")
	influxOrg.SetText(commons.GetInfluxOrg())
	influxToken := widget.NewEntry()
	influxToken.SetPlaceHolder("InfluxDB2 security token")
	influxToken.SetText(commons.GetInfluxToken())

	influxEnable := widget.NewCheck("Enable InfluxDB", func(onOff bool) {
		if onOff {
			commons.SetEnableInfluxDB("true")
		} else {
			commons.SetEnableInfluxDB("false")
		}
	})
	influxEnable.SetChecked(commons.IsInfluxDBEnabled())

	debugEnable := widget.NewCheck("Enable Debug Mode", func(onOff bool) {
		if onOff {
			commons.SetEnableDebugMode("true")
		} else {
			commons.SetEnableDebugMode("false")
		}
	})
	debugEnable.SetChecked(commons.IsDebugMode())

	testEnable := widget.NewCheck("Enable Test Mode", func(onOff bool) {
		if onOff {
			commons.SetEnableTestMode("true")
		} else {
			commons.SetEnableTestMode("false")
		}
	})
	testEnable.SetChecked(commons.IsTestMode())

	mqttUri := widget.NewEntry()
	mqttUri.SetPlaceHolder("MQTT ip tcp://10.100.1.16:1883")
	mqttUri.SetText(commons.GetMqttHostUri())
	mqttUser := widget.NewEntry()
	mqttUser.SetPlaceHolder("MQTT user name")
	mqttUser.SetText(commons.GetMqttUser())
	mqttPass := widget.NewPasswordEntry()
	mqttPass.SetPlaceHolder("MQTT user password")
	mqttPass.SetText(commons.GetMqttPass())

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "InfluxDB2 Host Url", Widget: influxUrl},
			{Text: "InfluxDB2 Bucket Name", Widget: influxBucket},
			{Text: "InfluxDB2 Organization", Widget: influxOrg},
			{Text: "InfluxDB2 Security Token", Widget: influxToken},

			{Text: "MQTT Host Url", Widget: mqttUri},
			{Text: "MQTT Username", Widget: mqttUser},
			{Text: "MQTT Password", Widget: mqttPass},

			{Text: "InfluxDB2 Enabled", Widget: influxEnable},
			{Text: "Debug Mode Enabled", Widget: debugEnable},
			{Text: "Test Mode Enabled", Widget: testEnable},
		},
		SubmitText: "Apply",
	}
	form.OnSubmit = func() { // optional, handle form submission
		commons.SetInfluxHostUri(influxUrl.Text)
		commons.SetInfluxBucket(influxBucket.Text)
		commons.SetInfluxOrg(influxOrg.Text)
		commons.SetInfluxToken(influxToken.Text)
		commons.SetMqttHostUri(mqttUri.Text)
		commons.SetMqttUser(mqttUser.Text)
		commons.SetMqttPass(mqttPass.Text)
		if fyne.CurrentApp() != nil {
			for key, value := range commons.GetConfigurationMap() {
				fyne.CurrentApp().Preferences().SetString(key, value)
			}
		}
		fmt.Println("Form submitted: restart for effect")
	}

	page := container.NewBorder(desc, nil, nil, nil, form)
	return page
}

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func SknTrayMenu(a fyne.App, w fyne.Window, chart fyne.Window) {
	// Add SystemBar Menu
	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("Mqtt2InfluxDB2",
			fyne.NewMenuItem("Show main", func() {
				w.Show()
			}),
			fyne.NewMenuItem("Show chart", func() {
				chart.Show()
			}))
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(theme.VisibilityIcon())

		w.SetCloseIntercept(func() { w.Hide() })
	}
}
func sknMenus(a fyne.App, w fyne.Window) {

	settingsItem := fyne.NewMenuItem("Settings", func() {
		w := a.NewWindow("Settings")
		page := container.NewVBox(
			settings.NewSettings().LoadAppearanceScreen(w),
			PreferencesPage(),
		)
		w.SetContent(page)
		w.Resize(fyne.NewSize(800, 600))
		w.Show()
	})

	cutItem := fyne.NewMenuItem("Cut", func() {
		shortcutFocused(&fyne.ShortcutCut{
			Clipboard: w.Clipboard(),
		}, w)
	})
	copyItem := fyne.NewMenuItem("Copy", func() {
		shortcutFocused(&fyne.ShortcutCopy{
			Clipboard: w.Clipboard(),
		}, w)
	})
	pasteItem := fyne.NewMenuItem("Paste", func() {
		shortcutFocused(&fyne.ShortcutPaste{
			Clipboard: w.Clipboard(),
		}, w)
	})

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://developer.fyne.io")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItem("Support", func() {
			u, _ := url.Parse("https://fyne.io/support/")
			_ = a.OpenURL(u)
		}),
	)
	file := fyne.NewMenu("File")
	if !fyne.CurrentDevice().IsMobile() {
		file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	}
	mainMenu := fyne.NewMainMenu(
		// a quit item will be appended to our first menu
		file,
		fyne.NewMenu("Edit", cutItem, copyItem, pasteItem),
		helpMenu,
	)
	w.SetMainMenu(mainMenu)
	w.SetMaster()
}
