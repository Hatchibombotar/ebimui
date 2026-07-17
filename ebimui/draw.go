package ebimui

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func FillSolid(screen *ebiten.Image, widget GenericWidget, fillColor color.Color) {
	vector.DrawFilledRect(
		screen,
		float32(widget.GetResultX()), float32(widget.GetResultY()),
		float32(widget.GetResultWidth()), float32(widget.GetResultHeight()),
		fillColor,
		false,
	)
}

// written by ai
func FillNineSlice(screen *ebiten.Image, widget GenericWidget, nineSliceImage *ebiten.Image, Top, Right, Bottom, Left int) {
	destX := widget.GetResultX()
	destY := widget.GetResultY()
	destW := widget.GetResultWidth()
	destH := widget.GetResultHeight()

	imgW, imgH := nineSliceImage.Bounds().Dx(), nineSliceImage.Bounds().Dy()

	// Destination rectangles
	midW := destW - Left - Right
	midH := destH - Top - Bottom

	if midW < 0 || midH < 0 {
		// Not enough space to draw nine-slice
		return
	}

	// Helper to draw a region
	drawRegion := func(sx, sy, sw, sh int, dx, dy, dw, dh float64) {
		src := nineSliceImage.SubImage(image.Rect(sx, sy, sx+sw, sy+sh)).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(dw/float64(sw), dh/float64(sh))
		op.GeoM.Translate(dx, dy)
		screen.DrawImage(src, op)
	}

	// Top-left
	drawRegion(0, 0, Left, Top, float64(destX), float64(destY), float64(Left), float64(Top))
	// Top
	drawRegion(Left, 0, imgW-Left-Right, Top, float64(destX+Left), float64(destY), float64(midW), float64(Top))
	// Top-right
	drawRegion(imgW-Right, 0, Right, Top, float64(destX+Left+midW), float64(destY), float64(Right), float64(Top))

	// Left
	drawRegion(0, Top, Left, imgH-Top-Bottom, float64(destX), float64(destY+Top), float64(Left), float64(midH))
	// Center
	drawRegion(Left, Top, imgW-Left-Right, imgH-Top-Bottom, float64(destX+Left), float64(destY+Top), float64(midW), float64(midH))
	// Right
	drawRegion(imgW-Right, Top, Right, imgH-Top-Bottom, float64(destX+Left+midW), float64(destY+Top), float64(Right), float64(midH))

	// Bottom-left
	drawRegion(0, imgH-Bottom, Left, Bottom, float64(destX), float64(destY+Top+midH), float64(Left), float64(Bottom))
	// Bottom
	drawRegion(Left, imgH-Bottom, imgW-Left-Right, Bottom, float64(destX+Left), float64(destY+Top+midH), float64(midW), float64(Bottom))
	// Bottom-right
	drawRegion(imgW-Right, imgH-Bottom, Right, Bottom, float64(destX+Left+midW), float64(destY+Top+midH), float64(Right), float64(Bottom))
}

func CreateFillSolid(fillColor color.Color) func(screen *ebiten.Image, widget GenericWidget, _ *UIContext) {
	return func(screen *ebiten.Image, widget GenericWidget, _ *UIContext) {
		FillSolid(screen, widget, fillColor)
	}
}

// TODO: add image helpers
