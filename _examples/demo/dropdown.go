package main

import (
	"image/color"

	"github.com/hatchibombotar/ebimui/ebimui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type DropdownOption struct {
	label string
	data  any
}

func (g *Game) CreateDropdown(options []string, selectedIndex *int) *ebimui.Box {
	ctx := g.uiContext
	dropdownId := ctx.GetId()

	dropdown := ctx.NewBoxWidget(func(b *ebimui.Box) {
		b.LayoutDirection(ebimui.LayoutColumn)
		b.Gap(2)

		b.AppendNewBoxWidget(func(b *ebimui.Box) {
			b.Padding(12+4, 12+8, 12+4, 12+8)
			if g.uiContext.IsWidgetHovered(b) {
				if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
					b.DrawNineSlice(input_outline_rectangle, 12)
				} else {
					b.DrawNineSlice(input_outline_rectangle, 12)
				}
			} else {
				b.DrawNineSlice(input_outline_rectangle, 12)
			}
			b.CursorShape(ebiten.CursorShapePointer)
			if ctx.IsWidgetHovered(b) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
				if g.expandedWidget != nil && *g.expandedWidget == dropdownId {
					g.expandedWidget = nil
				} else {
					g.expandedWidget = &dropdownId
				}
			}
			b.AppendNewTextWidget(func(t *ebimui.Text) {
				if *selectedIndex >= 0 && *selectedIndex < len(options) {
					t.Content(options[*selectedIndex])
				} else {
					t.Content("Select option...")
				}
				t.Face(loadedFonts.buttonFace)
				t.Color(color.RGBA{65, 65, 65, 255})
			})
		})

		// Display dropdown options only if expanded, positioned relative to not affect layout
		isExpanded := g.expandedWidget != nil && *g.expandedWidget == dropdownId
		if isExpanded {
			b.AppendNewBoxWidget(func(optionsContainer *ebimui.Box) {
				optionsContainer.LayoutDirection(ebimui.LayoutColumn)
				optionsContainer.WidthGrow()
				optionsContainer.PositionRelative(0, 57)

				for i, option := range options {
					optionsContainer.AppendNewBoxWidget(func(item *ebimui.Box) {
						item.Padding(12, 8, 12, 8)
						item.WidthGrow()
						item.CursorShape(ebiten.CursorShapePointer)

						if ctx.IsWidgetHovered(item) {
							item.DrawFillSolid(color.RGBA{235, 235, 235, 255})
							if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
								*selectedIndex = i
								g.expandedWidget = nil // Close after selection
							}
						} else {
							item.DrawFillSolid(color.RGBA{205, 205, 205, 255})
						}

						item.AlignHorizontal(ebimui.AlignCenter)

						item.AppendNewTextWidget(func(t *ebimui.Text) {
							t.Content(option)
							t.Face(loadedFonts.buttonFace)
							t.Color(color.RGBA{65, 65, 65, 255})
						})
					})
				}
			})
		}
	})
	return dropdown
}
