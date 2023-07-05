package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"net/url"
)

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}

func SknTrayMenu(a fyne.App, w fyne.Window) {
	// Add SystemBar Menu
	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("Mqtt2InfluxDB2",
			fyne.NewMenuItem("Show", func() {
				w.Show()
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
			layout.NewSpacer(),
		)
		w.SetContent(page)
		w.Resize(fyne.NewSize(640, 480))
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
