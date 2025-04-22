# Canvas : Circle

The Circle canvas object draws a styled circle in the center of the available space.
Its diameter will be the smaller of the available width and height, so it remains circular whilst filling the space.

## Basic usage

A circle is created using the `NewCircle` constructor function, passing in the color it should use to draw.
A basic circle is a simple fill, though it can also have a stroke (border) as you will see below.

```
circle := canvas.NewCircle(color.NRGBA{0, 0, 0x80, 0xff})
```

It is recommended to use the NRGBA color type to avoid complications with transparency.

## Setting styles

By default, a circle has no stroke (border), but you can set it up using the
`StrokeWidth` and `StrokeColor` fields.

```
circle := canvas.NewCircle(color.NRGBA{0x30, 0x30, 0x30, 0xb0})
circle.StrokeColor = color.NRGBA{0x80, 0, 0, 0xff}
circle.StrokeWidth = 2
```
