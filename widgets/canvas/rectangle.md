# Canvas : Line

The Rectang;e canvas object draws a styled rectangle that fills available space.
A basic rectangle is a simple fill, though it can have a stroke and also rounded corners, as you will see below.

## Basic usage

A rectangle is created using the `NewRectangle` constructor function, passing in the color it should use to draw for fill.
By default, a rectangle will have no stroke (border) and just a solit fill:

```
rect := canvas.NewRectangle(color.NRGBA{0x80, 0, 0, 0xff})
```

It is recommended to use the NRGBA color type to avoid complications with transparency.

## Setting styles

By default, a rectangle has no stroke (border), but you can set it up using the
`StrokeWidth` and `StrokeColor` fields.

```
rect := canvas.NewRectangle(color.NRGBA{0x80, 0, 0, 0xff})
rect.StrokeColor = color.NRGBA{0xff, 0x80, 0, 0xff}
rect.StrokeWidth = 3
```

## Rounded Corners

A rectangle (and its stroke) can also have rounded corners, controlled by the `CornerRadius` field:

```
rect := canvas.NewRectangle(color.NRGBA{0x80, 0, 0, 0xff})
rect.StrokeColor = color.NRGBA{0xff, 0, 0x80, 0xff}
rect.StrokeWidth = 2.0
rect.CornerRadius = 5
```