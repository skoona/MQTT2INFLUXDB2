package commons

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

func SknSelectThemedImage(alias string) *canvas.Image {
	return SknImageByName(alias, true)
}

func SknImageByName(alias string, themed bool) *canvas.Image {
	var selected fyne.Resource

	switch alias {
	case "garageOpen":
		selected = resourceGarageOpenSvg
	case "garageClosed":
		selected = resourceGarageClosedSvg
	default:
		selected = resourceSensorsOnMbo24pxSvg
	}
	image := canvas.NewImageFromResource(selected)
	if themed {
		image.Resource = theme.NewThemedResource(selected)
	}
	image.FillMode = canvas.ImageFillContain
	image.ScaleMode = canvas.ImageScaleSmooth
	return image
}
