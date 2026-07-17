package main

import (
	"fmt"
	"image/color"

	"github.com/hatchibombotar/ebimui/ebimui"

	"github.com/hajimehoshi/ebiten/v2"
)

// Slider:
// width = 300
// Bar at 0,0.

func (g *Game) CreateSlider(min, max int, value *int) *ebimui.Box {
	ctx := g.uiContext
	trackWidth := 300

	slider := ctx.NewBoxWidget(func(b *ebimui.Box) {
		b.LayoutDirection(ebimui.LayoutColumn)
		b.Gap(8)
		b.WidthGrow()

		// Display current value
		b.AppendNewTextWidget(func(t *ebimui.Text) {
			t.Content(fmt.Sprintf("Value: %d / %d", *value, max))
			t.Face(loadedFonts.buttonFace)
			t.Color(color.Black)
		})

		// Slider track container
		b.AppendNewBoxWidget(func(trackContainer *ebimui.Box) {
			trackContainer.FixedWidth(trackWidth)
			trackContainer.FixedHeight(24)
			trackContainer.Padding(8, 4, 8, 4)
			trackContainer.AlignVertical(ebimui.AlignCenter)
			trackContainer.CursorShape(ebiten.CursorShapePointer)

			// Handle track clicks
			trackContainer.DeferToPostLayout(func() {
				if ctx.IsWidgetHovered(trackContainer) && ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
					mx, _ := ebiten.CursorPosition()
					trackX := trackContainer.GetResultX()
					trackW := trackContainer.GetResultWidth()
					relativeX := mx - trackX
					if relativeX < 0 {
						relativeX = 0
					}
					if relativeX > trackW {
						relativeX = trackW
					}
					progress := float64(relativeX) / float64(trackW)
					fmt.Println(min + int(progress*float64(max-min)))
					*value = min + int(progress*float64(max-min))
				}
			})

			// Track background
			trackContainer.AppendNewBoxWidget(func(track *ebimui.Box) {
				track.FixedHeight(4)
				track.WidthGrow()
				track.DrawFillSolid(color.RGBA{180, 180, 180, 255})
			})

			// Thumb handle
			trackContainer.AppendNewBoxWidget(func(thumb *ebimui.Box) {
				thumb.FixedWidth(16)
				thumb.FixedHeight(16)

				// Calculate thumb position based on value
				range_ := max - min
				progress := float64(*value-min) / float64(range_)
				thumbX := int(progress * float64(trackWidth))
				thumb.PositionRelative(thumbX-8, 0)

				// Styling
				if ctx.IsWidgetHovered(thumb) || ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
					thumb.DrawFillSolid(color.RGBA{60, 120, 200, 255})
				} else {
					thumb.DrawFillSolid(color.RGBA{100, 140, 200, 255})
				}
				thumb.CursorShape(ebiten.CursorShapePointer)

				// Handle dragging
				b.DeferToPostLayout(func() {
					if ctx.IsWidgetHovered(thumb) && ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
						mx, _ := ebiten.CursorPosition()
						trackX := trackContainer.GetResultX() + 8
						trackW := trackWidth - 32

						relativeX := mx - trackX
						if relativeX < 0 {
							relativeX = 0
						}
						if relativeX > trackW {
							relativeX = trackW
						}

						progress := float64(relativeX) / float64(trackW)
						*value = min + int(progress*float64(max-min))
					}
				})
			})
		})
	})
	return slider
}
