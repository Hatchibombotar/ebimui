package main

import (
	"fmt"
	"image/color"

	"github.com/hatchibombotar/ebimui/ebimui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) CreateSlider(minValue, maxValue int, value *int) *ebimui.Box {
	ctx := g.uiContext
	trackWidth := 300
	scrubberWidth := 16

	slider := ctx.NewBoxWidget(func(b *ebimui.Box) {
		b.LayoutDirection(ebimui.LayoutColumn)
		b.Gap(8)
		b.WidthGrow()

		// Display current value
		b.AppendNewTextWidget(func(t *ebimui.Text) {
			t.Content(fmt.Sprintf("Value: %d / %d", *value, maxValue))
			t.Face(loadedFonts.buttonFace)
			t.Color(color.Black)
		})

		// Slider track container
		b.AppendNewBoxWidget(func(trackContainer *ebimui.Box) {
			trackContainer.FixedWidth(trackWidth)
			trackContainer.FixedHeight(24)
			trackContainer.Padding(8, 0, 8, 0)
			trackContainer.AlignVertical(ebimui.AlignCenter)
			trackContainer.CursorShape(ebiten.CursorShapePointer)

			isHot := ctx.IsWidgetHot(trackContainer)
			// Handle track clicks
			trackContainer.DeferToPostLayout(func() {
				if ctx.IsWidgetHovered(trackContainer) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
					ctx.SetWidgetHot(trackContainer)
				}

				if isHot {
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
					progress := (float64(relativeX) - float64(scrubberWidth/2)) / float64(trackW-scrubberWidth)
					progress = max(0, progress)
					progress = min(1, progress)
					*value = minValue + int(progress*float64(maxValue-minValue))
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
				range_ := maxValue - minValue
				progress := float64(*value-minValue) / float64(range_)
				thumbX := int(progress*float64(trackWidth-scrubberWidth)) + (scrubberWidth / 2)
				thumb.PositionRelative(thumbX-(scrubberWidth/2), 4)

				// Styling
				if isHot {
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
						*value = minValue + int(progress*float64(maxValue-minValue))
					}
				})
			})
		})
	})
	return slider
}
