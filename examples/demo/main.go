package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hatchibombotar/ebimui/ebimui"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	uiContext      *ebimui.UIContext
	expandedWidget *ebimui.WidgetID

	selectedTabIndex int

	selectedRadioItem    int
	selectedDropdownItem int
	checkboxTicked       bool

	sliderValue int

	gridItemCount int
}

var tabs = []string{"Button", "Checkbox", "Dropdown", "Grid Layout", "Radio", "Slider"}

// TODO: "Box", "Card Game"

func (g *Game) Update() error {
	g.uiContext.PreUpdate()

	g.uiContext.AddUi(
		g.uiContext.NewBoxWidget(func(b *ebimui.Box) {
			b.Padding(16, 16, 16, 16)
			b.LayoutDirection(ebimui.LayoutRow)
			b.FixedWidth(900)
			b.Gap(16)
			// Sidebar
			b.AppendNewBoxWidget(func(b *ebimui.Box) {
				b.Padding(12+4, 12+4, 12+4, 12+4)
				b.DrawNineSlice(button_rectangle_line, 12)
				for i, label := range tabs {
					b.AppendNewBoxWidget(func(b *ebimui.Box) {
						b.AppendNewTextWidget(func(t *ebimui.Text) {
							t.Content(label)
							t.Face(loadedFonts.tabFace)
							t.Color(color.Black)
						})
						b.CursorShape(ebiten.CursorShapePointer)

						if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) && g.uiContext.IsWidgetHovered(b) {
							g.uiContext.SetWidgetHot(b)
						}

						if g.uiContext.IsWidgetHovered(b) && g.uiContext.IsWidgetHot(b) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
							g.selectedTabIndex = i
						}
						if g.selectedTabIndex == i {
							b.DrawFillSolid(color.RGBA{205, 205, 205, 255})
						} else if g.uiContext.IsWidgetHovered(b) {
							b.DrawFillSolid(color.RGBA{235, 235, 235, 255})
						}
						b.Padding(12, 12, 12, 12)
						b.WidthGrow()
					})
				}
			})
			// Main Container
			b.AppendNewBoxWidget(func(b *ebimui.Box) {
				b.DrawNineSlice(box_nine_slice, 2)
				b.DrawNineSlice(button_rectangle_line, 12)
				b.WidthGrow()
				b.Padding(12+8, 12+8, 12+16, 12+8)
				b.Gap(16)

				b.AppendNewTextWidget(func(t *ebimui.Text) {
					t.Content(tabs[g.selectedTabIndex])
					t.Face(loadedFonts.titleFace)
					t.Color(color.Black)
				})

				switch tabs[g.selectedTabIndex] {
				case "Button":
					b.AppendNewBoxWidget(func(b *ebimui.Box) {
						b.Padding(12+4, 12+8, 12+4, 12+8)
						if g.uiContext.IsWidgetHovered(b) && inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
							g.uiContext.SetWidgetHot(b)
						}

						if g.uiContext.IsWidgetHot(b) {
							b.DrawNineSlice(button_rectangle_line_pressed, 12)
						} else if g.uiContext.IsWidgetHovered(b) {
							b.DrawNineSlice(button_rectangle_line_hover, 12)
						} else {
							b.DrawNineSlice(button_rectangle_line, 12)
						}

						b.CursorShape(ebiten.CursorShapePointer)
						if g.uiContext.IsWidgetHovered(b) && g.uiContext.IsWidgetHot(b) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
							fmt.Println("Press")
						}
						b.AppendNewTextWidget(func(t *ebimui.Text) {
							t.Content("Button")
							t.Face(loadedFonts.buttonFace)
							t.Color(color.RGBA{65, 65, 65, 255})
						})
					})

				case "Checkbox":
					// check_square_grey
					// check_square_grey_checkmark
					b.AppendNewBoxWidget(func(checkboxItem *ebimui.Box) {
						if g.uiContext.IsWidgetHovered(checkboxItem) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
							g.checkboxTicked = !g.checkboxTicked
						}

						checkboxItem.LayoutDirection(ebimui.LayoutRow)
						checkboxItem.Gap(6)
						checkboxItem.CursorShape(ebiten.CursorShapePointer)
						checkboxItem.AlignVertical(ebimui.AlignCenter)

						checkboxItem.AppendNewBoxWidget(func(b *ebimui.Box) {
							b.FixedWidth(32)
							b.FixedHeight(32)
							b.OnDraw(func(screen *ebiten.Image, widget ebimui.GenericWidget) {
								op := &ebiten.DrawImageOptions{}
								op.GeoM.Translate(float64(b.GetResultX()), float64(b.GetResultY()))
								if g.checkboxTicked {
									screen.DrawImage(check_square_grey_checkmark, op)
								} else {
									screen.DrawImage(check_square_grey, op)
								}
							})
						})

						checkboxItem.AppendNewTextWidget(func(t *ebimui.Text) {
							t.Content("Accept terms and conditions")
							t.Face(loadedFonts.buttonFace)
							t.Color(color.Black)
						})
					})
				case "Dropdown":
					// options := []DropdownOption{}
					b.AddChild(g.CreateDropdown([]string{"Option A", "Option B", "Option C", "Option D"}, &g.selectedDropdownItem))

				case "Grid Layout":
					b.AppendNewBoxWidget(func(b *ebimui.Box) {
						b.LayoutDirection(ebimui.LayoutRow)
						b.Gap(32)

						b.AppendNewBoxWidget(func(b *ebimui.Box) {
							b.LayoutMode(ebimui.LayoutGrid)
							b.Columns(3)
							b.Gap(4)
							n := g.gridItemCount
							for i := range n {
								b.AppendNewBoxWidget(func(b *ebimui.Box) {
									r := uint8(255.0 * (float64(n-i) / float64(n)))
									g := uint8(255.0 * ((float64(i)) / float64(n)))
									b.DrawFillSolid(color.CMYK{r, g, 50, 0})
									b.FixedHeight(40)
									b.FixedWidth(40)
								})
							}
						})

						b.AddChild(
							g.CreateSlider(3, 36, &g.gridItemCount),
						)
					})

					// var d int = 1
					// options := []DropdownOption{}
					// b.AddChild(g.CreateDropdown([]string{"Option A", "Option B", "Option C"}, &d))
				case "Radio":
					for i := range 5 {
						b.AppendNewBoxWidget(func(radioItem *ebimui.Box) {
							if g.uiContext.IsWidgetHovered(radioItem) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
								g.selectedRadioItem = i
							}

							radioItem.LayoutDirection(ebimui.LayoutRow)
							radioItem.Gap(4)
							radioItem.CursorShape(ebiten.CursorShapePointer)
							radioItem.AlignVertical(ebimui.AlignCenter)
							radioItem.AppendNewBoxWidget(func(b *ebimui.Box) {
								b.FixedWidth(32)
								b.FixedHeight(32)
								b.OnDraw(func(screen *ebiten.Image, widget ebimui.GenericWidget) {
									op := &ebiten.DrawImageOptions{}
									op.GeoM.Translate(float64(b.GetResultX()), float64(b.GetResultY()))
									if g.selectedRadioItem == i {
										screen.DrawImage(radio_selected, op)
									} else if g.uiContext.IsWidgetHovered(radioItem) {
										screen.DrawImage(radio_hovered, op)
									} else {
										screen.DrawImage(radio_deselected, op)
									}
								})
							})
							radioItem.AppendNewTextWidget(func(t *ebimui.Text) {
								t.Content(fmt.Sprint("Option ", i))
								t.Face(loadedFonts.buttonFace)
								t.Color(color.Black)
							})
						})
					}
				case "Slider":
					b.AppendNewBoxWidget(func(b *ebimui.Box) {
						b.LayoutDirection(ebimui.LayoutColumn)
						b.Gap(8)
						b.WidthGrow()

						// Display current value
						b.AppendNewTextWidget(func(t *ebimui.Text) {
							t.Content(fmt.Sprintf("Value: %d / %d", g.sliderValue, 100))
							t.Face(loadedFonts.buttonFace)
							t.Color(color.Black)
						})

						b.AddChild(g.CreateSlider(0, 100, &g.sliderValue))
					})
				}

			})
		}),
	)

	g.uiContext.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{39, 117, 166, 255})
	g.uiContext.DrawChildren(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 900, 675
}

func main() {
	ebiten.SetWindowSize(900*0.8, 675*0.8)
	ebiten.SetWindowTitle("ebimui demo")

	g := &Game{
		uiContext: ebimui.NewUIContext(),

		gridItemCount: 15,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
