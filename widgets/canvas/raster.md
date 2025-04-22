# Canvas : Raster

The Raster canvas object is used to draw image content to screen programmatically.
A raster object is used to fill pixels on screen, where the number of pixels will vary based on the
density of the output device and the space available.

## Basic usage

A raster is typically created using the `NewRaster` constructor which passes in a generator function.
The generator function is called when the raster refreshes or the number of pixels required changes.

```
rast := canvas.NewRaster(func(w, h int) image.Image {
	ret := image.NewNRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if (x+y)/4%2 == 0 {
				ret.Set(x, y, color.Black)
			} else {
				ret.Set(x, y, color.White)
			}
		}
	}
	return ret
})
```
