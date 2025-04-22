# Canvas : Image

The Image canvas object is used to draw a loaded image to the screen.
By default, the image will fill available space, be sure to place it in container where space is allocated appropriately.

## Basic usage

An image is typically created using the `NewImageFromFile`, `NewImgeFromResource`, or `NewImageFromImage` constructor functions,
passing the appropriate parameter to specify the location of the image to display.

```
img := canvas.NewImageFromResource(theme.HomeIcon())
```

Note that the default minimum size for an image is just 1x1 because a Fyne app does not use pixel sizes for layout.

## Setting fill and scale modes

You can control scale and fill modes to adjust how your iamge fills the space.
It is also possible to make the image blend into the background using `Translucency`.

```
img := canvas.NewImageFromResource(theme.ComputerIcon())
img.FillMode = canvas.ImageFillContain
img.ScaleMode = canvas.ImageScaleFastest
img.Translucency = 0.5
```
