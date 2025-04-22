# Canvas : Text

The Text canvas object is used to draw text to screen with control over its output.

## Basic usage

Text is created using the `NewText` constructor function, passing in the string to
display along with the color it should use to draw.

```
txt := canvas.NewText("Hello", color.NRGBA{0, 0x80, 0, 0xff})
```

It is recommended to use the NRGBA color type to avoid complications with transparency.

## Setting styles

By default, text is the normal style (not italic or bold or monospace) and
the font size matches the operating system normal text (approx 14pt).

You can specify different values for these and other fields as shown:

```
txt := canvas.NewText("Styled", color.NRGBA{0x80, 0, 0, 0xb0})
txt.TextSize = 24
txt.TextStyle.Bold = true
txt.Alignment = fyne.TextAlignTrailing
```
