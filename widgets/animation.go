package widgets

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func makeAnimationScreen(_ fyne.Window, Data interface{}) fyne.CanvasObject {
	curves := makeAnimationCurves()
	curves.Move(fyne.NewPos(0, 140+theme.Padding()))
	return container.NewWithoutLayout(makeAnimationCanvas(), curves)
}

func makeAnimationCanvas() fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Black)
	rect.Resize(fyne.NewSize(410, 140))

	a := canvas.NewColorRGBAAnimation(
		color.NRGBA{R: 0x29, G: 0x6f, B: 0xf6, A: 0xaa},
		color.NRGBA{R: 0x8b, G: 0xc3, B: 0x4a, A: 0xaa},
		time.Second*3, func(c color.Color) {
			rect.FillColor = c
			canvas.Refresh(rect)
		},
	)
	a.RepeatCount = fyne.AnimationRepeatForever
	a.AutoReverse = true
	a.Start()

	var a2 *fyne.Animation
	i := widget.NewIcon(theme.CheckButtonCheckedIcon())
	a2 = canvas.NewPositionAnimation(fyne.NewPos(0, 0), fyne.NewPos(350, 80), time.Second*3, func(p fyne.Position) {
		i.Move(p)

		width := 10 + (p.X / 7)
		i.Resize(fyne.NewSize(width, width))
	})
	a2.RepeatCount = fyne.AnimationRepeatForever
	a2.AutoReverse = true
	a2.Curve = fyne.AnimationLinear
	a2.Start()

	OnChangeFuncs = append(OnChangeFuncs, a.Stop, a2.Stop)

	running := true
	var toggle *widget.Button
	toggle = widget.NewButton("Stop", func() {
		if running {
			a.Stop()
			a2.Stop()
			toggle.SetText("Start")
		} else {
			a.Start()
			a2.Start()
			toggle.SetText("Stop")
		}
		running = !running
	})
	toggle.Resize(toggle.MinSize())
	toggle.Move(fyne.NewPos(152, 54))
	return container.NewWithoutLayout(rect, i, toggle)
}

func makeAnimationCurves() fyne.CanvasObject {
	label1, box1, a1 := makeAnimationCurveItem("EaseInOut", fyne.AnimationEaseInOut, 0)
	label2, box2, a2 := makeAnimationCurveItem("EaseIn", fyne.AnimationEaseIn, 30+theme.Padding())
	label3, box3, a3 := makeAnimationCurveItem("EaseOut", fyne.AnimationEaseOut, 60+theme.Padding()*2)
	label4, box4, a4 := makeAnimationCurveItem("Linear", fyne.AnimationLinear, 90+theme.Padding()*3)

	OnChangeFuncs = append(OnChangeFuncs, a1.Stop, a2.Stop, a3.Stop, a4.Stop)

	start := widget.NewButton("Compare", func() {
		a1.Start()
		a2.Start()
		a3.Start()
		a4.Start()
	})
	start.Resize(start.MinSize())
	start.Move(fyne.NewPos(0, 120+theme.Padding()*4))
	return container.NewWithoutLayout(label1, label2, label3, label4, box1, box2, box3, box4, start)
}

func makeAnimationCurveItem(label string, curve fyne.AnimationCurve, yOff float32) (
	text *widget.Label, box fyne.CanvasObject, anim *fyne.Animation) {
	text = widget.NewLabel(label)
	text.Alignment = fyne.TextAlignCenter
	text.Resize(fyne.NewSize(380, 30))
	text.Move(fyne.NewPos(0, yOff))
	box = newThemedBox()
	box.Resize(fyne.NewSize(30, 30))
	box.Move(fyne.NewPos(0, yOff))

	anim = canvas.NewPositionAnimation(
		fyne.NewPos(0, yOff), fyne.NewPos(380, yOff), time.Second, func(p fyne.Position) {
			box.Move(p)
			box.Refresh()
		})
	anim.Curve = curve
	anim.AutoReverse = true
	anim.RepeatCount = 1
	return
}

// themedBox is a simple box that change its background color according
// to the selected theme
type themedBox struct {
	widget.BaseWidget
}

func newThemedBox() *themedBox {
	b := &themedBox{}
	b.ExtendBaseWidget(b)
	return b
}

func (b *themedBox) CreateRenderer() fyne.WidgetRenderer {
	b.ExtendBaseWidget(b)
	bg := canvas.NewRectangle(theme.Color(theme.ColorNameForeground))
	return &themedBoxRenderer{bg: bg, objects: []fyne.CanvasObject{bg}}
}

type themedBoxRenderer struct {
	bg      *canvas.Rectangle
	objects []fyne.CanvasObject
}

func (r *themedBoxRenderer) Destroy() {
}

func (r *themedBoxRenderer) Layout(size fyne.Size) {
	r.bg.Resize(size)
}

func (r *themedBoxRenderer) MinSize() fyne.Size {
	return r.bg.MinSize()
}

func (r *themedBoxRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *themedBoxRenderer) Refresh() {
	r.bg.FillColor = theme.Color(theme.ColorNameForeground)
	r.bg.Refresh()
}
