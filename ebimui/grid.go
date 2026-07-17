package ebimui

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// Init Functions
// -----------

func (uiContext *UIContext) NewGridWidget(options func(b *Grid)) *Grid {
	box := &Grid{
		op:  &GridWidgetOps{},
		id:  uiContext.GetId(),
		ctx: uiContext,
	}

	options(box)

	return box
}

// Rename to AddNewGridWidget?
func (box *Grid) AppendNewGridWidget(options func(b *Grid)) {
	ctx := box.ctx
	box.AddChild(
		ctx.NewGridWidget(options),
	)
}
func (box *Grid) AppendNewTextWidget(options func(t *Text)) {
	ctx := box.ctx
	box.AddChild(
		ctx.NewTextWidget(options),
	)
}
func (box *Grid) LayoutDirection(direction LayoutDirection) {
	box.op.LayoutDirection = direction
}

// TODO: rework how layout + columns work as columns only make sense for grid
func (box *Grid) LayoutMode(mode LayoutType) {
	box.op.LayoutType = mode
}
func (box *Grid) Columns(columns int) {
	box.op.Columns = columns
}

func (box *Grid) Padding(Top, Right, Bottom, Left int) {
	box.op.Padding = Spacing{Top, Right, Bottom, Left}
}

func (box *Grid) FixedWidth(Width int) {
	box.op.WidthMode = SizeFixed
	box.op.Width = Width
}

func (box *Grid) FixedHeight(Height int) {
	box.op.HeightMode = SizeFixed
	box.op.Height = Height
}

func (box *Grid) WidthGrow() {
	box.op.WidthMode = SizeGrow
}

func (box *Grid) Gap(Gap int) {
	box.op.Gap = Gap
}

func (box *Grid) AlignHorizontal(align AlignMode) {
	box.op.AlignHorizontal = align
}
func (box *Grid) AlignVertical(align AlignMode) {
	box.op.AlignVertical = align
}
func (box *Grid) PositionFixed(x, y int) {
	box.op.PositionMode = PositionFixed
	box.op.X = x
	box.op.Y = y
}
func (box *Grid) PositionRelative(x, y int) {
	box.op.PositionMode = PositionRelative
	box.op.X = x
	box.op.Y = y
}

// Consider removing root
func (box *Grid) OnDraw(drawFunc func(screen *ebiten.Image, widget GenericWidget, root *UIContext)) {
	box.op.OnDraw = drawFunc
}
func (box *Grid) DrawNineSlice(nineSliceImage *ebiten.Image, nineSliceWidth int) {
	box.op.OnDraw = func(screen *ebiten.Image, widget GenericWidget, root *UIContext) {
		FillNineSlice(screen, widget, nineSliceImage, nineSliceWidth, nineSliceWidth, nineSliceWidth, nineSliceWidth)
	}
}
func (box *Grid) DrawNineSlice4(nineSliceImage *ebiten.Image, Top, Right, Bottom, Left int) {
	box.op.OnDraw = func(screen *ebiten.Image, widget GenericWidget, root *UIContext) {
		FillNineSlice(screen, widget, nineSliceImage, Top, Right, Bottom, Left)
	}
}

func (box *Grid) DrawFillSolid(color color.Color) {
	box.op.OnDraw = func(screen *ebiten.Image, widget GenericWidget, root *UIContext) {
		FillSolid(screen, widget, color)
	}
}

func (box *Grid) CursorShape(cursorShape ebiten.CursorShapeType) {
	box.op.CursorShape = cursorShape
}

// --------
// Grid

type Grid struct {
	op *GridWidgetOps

	Children []GenericWidget
	// parent   GenericWidget

	eventHandler *UIContext

	resultX, resultY          int
	resultWidth, resultHeight int
	minWidth, minHeight       int

	id WidgetID

	ctx *UIContext
}

func (b *Grid) GetUniqueId() WidgetID {
	return b.id
}

type GridWidgetOps struct {
	WidthMode  SizeMode
	Width      int
	HeightMode SizeMode
	Height     int

	Padding Spacing

	// DrawBackgroundColour color.Color
	Draw func(w Grid, screen *ebiten.Image)

	LayoutType      LayoutType
	LayoutDirection LayoutDirection
	Gap             int

	// Rows or column count for LayoutType == Grid
	Rows, Columns int

	// align along the horizonal and vertical axies
	AlignHorizontal, AlignVertical AlignMode

	OnDraw func(screen *ebiten.Image, widget GenericWidget, root *UIContext)

	CursorShape ebiten.CursorShapeType

	PositionMode PositionMode

	// X and Y position of widget
	//
	// If Op.PositionMode is set to PositionFixed, relative to 0, 0
	//
	// If Op.PositionMode is set to PositionRelative, relative to Top Left & Right of Parent
	//
	// Otherwise, no other effects
	X, Y int
}

func (w *Grid) AddChild(child GenericWidget) {
	w.Children = append(w.Children, child)
}

func (w *Grid) Draw(screen *ebiten.Image, root *UIContext) {
	if w.op.OnDraw != nil {
		w.op.OnDraw(screen, w, root)
	}

	for _, child := range w.Children {
		child.Draw(screen, root)
	}
}

// TODO: Rename to LayoutUpdate*** ?
func (w *Grid) Update() {
	w.UpdateSizeWidthFitPass()
	w.UpdateSizeWidthGrowPass()
	w.UpdateSizeWrapWidth() // Used in text to wrap
	w.UpdateSizeHeightFitPass()
	w.UpdateSizeHeightGrowPass()
	w.UpdatePosition(0, 0)
}

func (w *Grid) UpdateSizeWidthFitPass() {
	if w.op.HeightMode == SizeFixed {
		w.resultWidth = w.op.Width
		w.minWidth = w.op.Width
		for _, child := range w.Children {
			child.UpdateSizeWidthFitPass()
		}
		return
	}
	if w.op.LayoutType == LayoutGrid {
		childrenPlaced := 0
		width := 0
		minWidth := 0
		rows := int(math.Ceil(float64(len(w.Children)) / float64(w.op.Columns)))
		for range rows {
			rowTotalWidth := 0
			rowTotalMinWidth := 0

			for range w.op.Columns {
				if childrenPlaced >= len(w.Children) {
					break
				}

				child := w.Children[childrenPlaced]
				child.UpdateSizeWidthFitPass()
				// TODO: Account for fixed & relative position mode for children

				rowTotalWidth += child.GetResultWidth()
				rowTotalMinWidth += child.GetMinWidth()
				childrenPlaced += 1
			}
			if rowTotalWidth > width {
				width = rowTotalWidth
			}
			if rowTotalMinWidth > minWidth {
				minWidth = rowTotalMinWidth
			}
		}
		width += w.op.Padding.Left + w.op.Padding.Right + (w.op.Columns-1)*w.op.Gap
		minWidth += w.op.Padding.Left + w.op.Padding.Right + (w.op.Columns-1)*w.op.Gap

		switch w.op.WidthMode {
		case SizeFit, SizeGrow:
			w.resultWidth = width
			w.minWidth = minWidth
		case SizeFixed:
			w.resultWidth = w.op.Width
			w.minWidth = w.op.Width
		default:
			panic("Unexpected sizing mode")
		}

	} else if w.op.LayoutDirection == LayoutRow {
		minWidth := w.op.Padding.Left + w.op.Padding.Right + (len(w.Children)-1)*w.op.Gap
		width := minWidth

		for _, child := range w.Children {
			child.UpdateSizeWidthFitPass()
			width += child.GetResultWidth()
			minWidth += child.GetMinWidth()
		}

		switch w.op.WidthMode {
		case SizeFit, SizeGrow:
			w.resultWidth = width
			w.minWidth = minWidth
		case SizeFixed:
			w.resultWidth = w.op.Width
			w.minWidth = w.op.Width
		default:
			panic("Unexpected sizing mode")
		}

	} else if w.op.LayoutDirection == LayoutColumn {
		width := 0
		minWidth := 0

		for _, child := range w.Children {
			child.UpdateSizeWidthFitPass()

			width = max(width, child.GetResultWidth())
			minWidth = max(minWidth, child.GetMinWidth())
		}
		width += w.op.Padding.Left + w.op.Padding.Right
		minWidth += w.op.Padding.Left + w.op.Padding.Right

		switch w.op.WidthMode {
		case SizeFit, SizeGrow:
			w.resultWidth = width
			w.minWidth = minWidth
		case SizeFixed:
			w.resultWidth = w.op.Width
			w.minWidth = w.op.Width
		default:
			panic("Unexpected sizing mode")
		}
	}
}

func (w *Grid) UpdateSizeHeightFitPass() {
	if w.op.HeightMode == SizeFixed {
		w.resultHeight = w.op.Height
		w.minHeight = w.op.Height
		for _, child := range w.Children {
			child.UpdateSizeHeightFitPass()
		}
		return
	}
	if w.op.LayoutType == LayoutGrid {
		childrenPlaced := 0
		height := 0
		minHeight := 0
		rows := int(math.Ceil(float64(len(w.Children)) / float64(w.op.Columns)))
		for range rows {
			rowMaxHeight := 0
			rowMaxMinHeight := 0

			for range w.op.Columns {
				if childrenPlaced >= len(w.Children) {
					break
				}

				child := w.Children[childrenPlaced]
				child.UpdateSizeHeightFitPass()
				// TODO: Account for fixed & relative position mode for children

				rowMaxHeight = max(rowMaxHeight, child.GetResultHeight())
				rowMaxMinHeight = max(rowMaxMinHeight, child.GetMinHeight())
				childrenPlaced += 1
			}
			height += rowMaxHeight
			minHeight += rowMaxMinHeight
		}
		height += w.op.Padding.Top + w.op.Padding.Bottom + (rows-1)*w.op.Gap
		minHeight += w.op.Padding.Top + w.op.Padding.Bottom + (rows-1)*w.op.Gap

		w.resultHeight = height
		w.minHeight = minHeight

	}
}

// TODO: Error when things could grow forever
func (parent *Grid) UpdateSizeWidthGrowPass() {
	remainingWidth := parent.GetResultWidth()

	if parent.op.LayoutType == LayoutFlex {
		if parent.op.LayoutDirection == LayoutColumn {
			maxWidth := remainingWidth - (parent.op.Padding.Left + parent.op.Padding.Right)
			for _, child := range parent.Children {
				if child.CanGrowWidth() {
					child.SetResultWidth(maxWidth)
				}
			}

			for _, child := range parent.Children {
				child.UpdateSizeWidthGrowPass()
			}
			return
		}
	}
	if parent.op.LayoutType == LayoutGrid {
		if parent.op.LayoutDirection == LayoutColumn {
			if parent.op.Columns == 0 {
				panic("Columns set to 0.")
			}
			remainingWidth /= parent.op.Columns
			remainingWidth -= parent.op.Gap * (parent.op.Columns - 1)
		} else {
			panic("Unexpected Layout Direction")
		}
	}
	remainingWidth -= parent.op.Padding.Left + parent.op.Padding.Right

	growable := make([]GenericWidget, 0)
	shrinkable := make([]GenericWidget, 0)
	for _, child := range parent.Children {
		remainingWidth -= child.GetResultWidth()
		if child.CanGrowWidth() {
			growable = append(growable, child)
		}
		if child.GetMinWidth() != child.GetResultWidth() {
			shrinkable = append(shrinkable, child)
		}
	}

	remainingWidth -= (len(parent.Children) - 1) * parent.op.Gap

	// grow elements
	for remainingWidth > 0 {
		if len(growable) == 0 {
			break
		}
		smallest := growable[0].GetResultWidth()
		secondSmallest := 99999999
		widthToAdd := remainingWidth
		for _, child := range growable {
			if child.GetResultWidth() < smallest {
				secondSmallest = smallest
				smallest = child.GetResultWidth()
			}
			if child.GetResultWidth() > smallest {
				secondSmallest = min(secondSmallest, child.GetResultWidth())
				widthToAdd = secondSmallest - smallest
			}
		}

		widthToAdd = min(widthToAdd, remainingWidth/len(growable))

		// temp
		if widthToAdd == 0 {
			widthToAdd = remainingWidth
		}

		for _, child := range growable {
			if child.GetResultWidth() == smallest {
				child.SetResultWidth(child.GetResultWidth() + widthToAdd)
				remainingWidth -= widthToAdd
			}
		}
	}

	// shrink elements
	for remainingWidth < 0 {
		if len(shrinkable) == 0 {
			break
		}
		largest := shrinkable[0].GetResultWidth()
		secondLargest := 99999999
		widthToAdd := remainingWidth
		for _, child := range shrinkable {
			if child.GetResultWidth() > largest {
				secondLargest = largest
				largest = child.GetResultWidth()
			}
			if child.GetResultWidth() < largest {
				secondLargest = max(secondLargest, child.GetResultWidth())
				widthToAdd = secondLargest - largest
			}
		}

		widthToAdd = min(widthToAdd, remainingWidth/len(shrinkable))

		// temp
		if widthToAdd == 0 {
			widthToAdd = remainingWidth
		}

		newShrinkable := make([]GenericWidget, 0)
		for _, child := range shrinkable {
			previousWidth := child.GetResultWidth()
			if child.GetResultWidth() == largest {
				child.SetResultWidth(child.GetResultWidth() + widthToAdd)
				if child.GetResultWidth() < child.GetMinWidth() {
					child.SetResultWidth(child.GetMinWidth())
					continue
				}
				remainingWidth -= child.GetResultWidth() - previousWidth
				if child.GetResultWidth() == child.GetMinWidth() {
					continue
				}
			}
			newShrinkable = append(newShrinkable, child)
		}
		shrinkable = newShrinkable
	}

	for _, child := range parent.Children {
		child.UpdateSizeWidthGrowPass()
	}
}

// TODO: this doesn't work, make it like width
func (parent *Grid) UpdateSizeHeightGrowPass() {
	remainingHeight := parent.GetResultHeight()
	remainingHeight -= parent.op.Padding.Top + parent.op.Padding.Bottom

	for _, child := range parent.Children {
		if !child.CanGrowHeight() {
			remainingHeight -= child.GetResultHeight()
		}
	}

	for _, child := range parent.Children {
		if child.CanGrowHeight() {
			child.SetResultHeight(remainingHeight)
		}
	}

	for _, child := range parent.Children {
		child.UpdateSizeHeightGrowPass()
	}
}

func (w *Grid) UpdateSizeWrapWidth() {
	for _, child := range w.Children {
		child.UpdateSizeWrapWidth()
	}
}

func (w *Grid) UpdatePosition(x, y int) {
	switch w.op.PositionMode {
	case PositionFixed:
		w.resultX = w.op.X
		w.resultY = w.op.Y
	case PositionRelative:
		w.resultX = x + w.op.X
		w.resultY = y + w.op.Y
	default:
		w.resultX = x
		w.resultY = y
	}

	if w.op.LayoutType == LayoutGrid {
		startX := w.resultX + w.op.Padding.Left
		currentY := w.resultY + w.op.Padding.Top
		childrenPlaced := 0
		rows := int(math.Ceil(float64(len(w.Children)) / float64(w.op.Columns)))
		for range rows {
			rowMaxHeight := 0

			currentX := startX
			for range w.op.Columns {
				if childrenPlaced >= len(w.Children) {
					break
				}

				child := w.Children[childrenPlaced]
				child.UpdatePosition(currentX, currentY)
				// TODO: Account for fixed & relative position mode for children

				rowMaxHeight = max(rowMaxHeight, child.GetResultHeight())
				childrenPlaced += 1

				currentX += child.GetResultWidth() + w.op.Gap

				// TODO: amount shifted should not be based on item width.
				// currentX += (w.resultWidth / w.Op.Columns) + w.Op.Gap
			}

			currentY += rowMaxHeight + w.op.Gap
		}
	}
}

// generic things

func (w *Grid) GetResultWidth() int {
	return w.resultWidth
}

func (w *Grid) GetResultHeight() int {
	return w.resultHeight
}

func (w *Grid) GetResultX() int {
	return w.resultX
}

func (w *Grid) GetResultY() int {
	return w.resultY
}

func (w *Grid) CanGrowWidth() bool {
	return w.op.WidthMode == SizeGrow
}
func (w *Grid) CanGrowHeight() bool {
	return w.op.HeightMode == SizeGrow
}

func (w *Grid) GetMinWidth() int {
	return w.minWidth
}
func (w *Grid) GetMinHeight() int {
	return w.minHeight
}

func (w *Grid) SetResultWidth(value int) {
	w.resultWidth = value
}
func (w *Grid) SetResultHeight(value int) {
	w.resultHeight = value
}

func (w *Grid) String() string {
	return w.string("")
}

func (w *Grid) string(padding string) string {
	str := padding + fmt.Sprint(
		"<Grid",
		" id=", w.GetUniqueId(),
		" width(", w.op.WidthMode, ")=", w.resultWidth,
		", height(", w.op.HeightMode, ")=", w.resultHeight,
		", padding=(", w.op.Padding.Left, ",", w.op.Padding.Top, ",", w.op.Padding.Right, ",", w.op.Padding.Bottom, ")",
		", x=", w.resultX,
		", y=", w.resultY,
		", hasEventHandler=", w.eventHandler != nil,
		">") + "\n"
	for _, child := range w.Children {
		switch child := child.(type) {
		case *Grid:
			str += child.string(padding+"  ") + "\n"
		default:
			str += padding + "  " + child.String() + "\n"
		}
	}
	str += padding + "</Grid>"
	return str
}

func (w *Grid) UpdateInput(ctx *UIContext) {
	ctx.elementCache[w.id] = w

	// if w.Op.IsFocusable && ctx.isWidgetHovered(w) && inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
	// 	root_DO_NOT_USE.SetFocusOn(w)
	// }

	if w.op.CursorShape != 0 && ctx.isWidgetHovered(w) {
		ctx.SetCursorShape(w.op.CursorShape)
	}

	if ctx.isWidgetHovered(w) {
		ctx.hoveredWidget = w.GetUniqueId()
		ctx.hoveredWidgets = append(ctx.hoveredWidgets, w.GetUniqueId())
	}

	for _, child := range w.Children {
		child.UpdateInput(ctx)
	}
}

func (w *Grid) GetPositionMode() PositionMode {
	return w.op.PositionMode
}
