package commons

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
)

func SknSelectResource(alias string) fyne.Resource {
	return SknImageByName(alias, false).Resource
}
func SknSelectThemedResource(alias string) fyne.Resource {
	return SknImageByName(alias, true).Resource
}

func SknSelectImage(alias string) *canvas.Image {
	return SknImageByName(alias, false)
}
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
	case "sensorOff_r":
		selected = resourceSensorsOffMbr24pxSvg
	case "sensorOn_o":
		selected = resourceSensorsOnMbo24pxSvg
	case "sensorOn_r":
		selected = resourceSensorsOnMbr24pxSvg
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
