package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/fyne-io/demo/data"
)

func rgbGradient(x, y, w, h int) color.Color {
	g := int(float32(x) / float32(w) * float32(255))
	b := int(float32(y) / float32(h) * float32(255))

	return color.NRGBA{uint8(255 - b), uint8(g), uint8(b), 0xff}
}

// k9Doc loads a graphics example panel for the demo app
func k9Doc(_ fyne.Window, Data interface{}) fyne.CanvasObject {
	return container.NewGridWrap(fyne.NewSize(90, 90),
		canvas.NewImageFromResource(data.FyneLogo),
	)
}
