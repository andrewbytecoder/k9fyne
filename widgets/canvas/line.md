# Canvas : Line

The Line canvas object draws a styled line from one position to another.
The default positions correspond with the top, left and bottom, right of its bounding box, 
though they can be controlled using Position1 and Position2 fields.
Therefor a line should normally be drawn in a container without layout, or a custom widget.

## Basic usage

A line is created using the `NewLine` constructor function, passing in the color it should use to draw.

```
line := canvas.NewLine(color.NRGBA{0, 0, 0x80, 0xff})
```

It is recommended to use the NRGBA color type to avoid complications with transparency.

## Setting styles

By default, a line is 1 unit across (device independent pixel).
You can specify different values using the `StrokeWidth` field.

```
line := canvas.NewLine(color.NRGBA{0, 0x80, 0, 0xdf})
line.StrokeWidth = 5
```
