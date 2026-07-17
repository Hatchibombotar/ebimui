package ebimui

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Box struct {
	op *BoxWidgetOps

	Children []GenericWidget
	// parent   GenericWidget

	eventHandler *UIContext

	resultX, resultY          int
	resultWidth, resultHeight int
	minWidth, minHeight       int

	id WidgetID

	ctx *UIContext

	// Options
	postLayoutFunction func()
}

func (b *Box) GetUniqueId() WidgetID {
	return b.id
}

type BoxWidgetOps struct {
	WidthMode  SizeMode
	Width      int
	HeightMode SizeMode
	Height     int

	Padding Spacing

	// DrawBackgroundColour color.Color
	Draw func(w Box, screen *ebiten.Image)

	LayoutType      LayoutType
	LayoutDirection LayoutDirection
	Gap             int

	// Rows or column count for LayoutType == Grid
	Rows, Columns int

	// align along the horizonal and vertical axies
	AlignHorizontal, AlignVertical AlignMode

	OnDraw func(screen *ebiten.Image, widget GenericWidget)

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

func (w *Box) AddChild(child GenericWidget) {
	w.Children = append(w.Children, child)
}

func (w *Box) Draw(screen *ebiten.Image, root *UIContext) {
	if w.op.OnDraw != nil {
		w.op.OnDraw(screen, w)
	}

	for _, child := range w.Children {
		child.Draw(screen, root)
	}
}

// TODO: Rename to LayoutUpdate*** ?
// Note: This is currently just ran once for each root node, could this cause issues?
func (w *Box) Update() {
	w.UpdateSizeWidthFitPass()
	// Post order depth first traversal
	// Parent widgets use `ResultWidth` and `MinWidth`

	w.UpdateSizeWidthGrowPass()
	// Pre order depth first traversal
	// Parent widgets set the `ResultWidth` of their children

	w.UpdateSizeWrapWidth()
	// Used in text to wrap

	w.UpdateSizeHeightFitPass()
	// Post order depth first traversal
	// Parent widgets use `ResultWidth` and `MinWidth`

	w.UpdateSizeHeightGrowPass()
	// Pre order depth first traversal
	// Parent widgets set the `ResultHeight` of their children

	w.UpdatePosition(0, 0)
	// Post order depth first traversal
	// Parent widgets are placed based on the position provided to the function
	// Child widgets are placed based on the `ResultWidth/ResultHeight` of sibilings

}

func (w *Box) UpdateSizeWidthFitPass() {
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

func (w *Box) UpdateSizeHeightFitPass() {
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

	} else if w.op.LayoutDirection == LayoutRow {
		minHeight := 0
		height := 0

		for _, child := range w.Children {
			child.UpdateSizeHeightFitPass()

			minHeight = max(minHeight, child.GetMinHeight())
			height = max(height, child.GetResultHeight())
		}
		height += w.op.Padding.Top + w.op.Padding.Bottom
		minHeight += w.op.Padding.Top + w.op.Padding.Bottom

		w.resultHeight = height
		w.minHeight = minHeight

	} else if w.op.LayoutDirection == LayoutColumn {
		minHeight := w.op.Padding.Top + w.op.Padding.Bottom + (len(w.Children)-1)*w.op.Gap
		height := minHeight

		for _, child := range w.Children {
			child.UpdateSizeHeightFitPass()

			switch child.GetPositionMode() {
			case PositionRelative, PositionFixed:
				minHeight -= w.op.Gap
				height -= w.op.Gap
				continue
			}

			height += child.GetResultHeight()
			minHeight += child.GetMinHeight()
		}

		w.resultHeight = height
		w.minHeight = minHeight
	}
}

// TODO: Should error when things could grow forever?
func (parent *Box) UpdateSizeWidthGrowPass() {
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

	// Categorise children into growable and shinkable widgets:

	// Elements with width set to grow
	growable := make([]GenericWidget, 0)
	// Elements larget than their minimum size
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
func (parent *Box) UpdateSizeHeightGrowPass() {
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

func (w *Box) UpdateSizeWrapWidth() {
	for _, child := range w.Children {
		child.UpdateSizeWrapWidth()
	}
}

func (w *Box) UpdatePosition(x, y int) {
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
	} else if w.op.LayoutDirection == LayoutRow {
		currentX := 0
		if w.op.AlignHorizontal == AlignStart {
			currentX = w.resultX + w.op.Padding.Left
		} else {
			panic("not implemented")
		}

		for _, child := range w.Children {
			if child.GetPositionMode() == PositionFixed || child.GetPositionMode() == PositionRelative {
				child.UpdatePosition(w.resultX, w.resultY)
				continue
			}
			currentY := 0
			switch w.op.AlignVertical {
			case AlignStart:
				currentY = w.resultY + w.op.Padding.Top
			case AlignCenter:
				currentY = w.resultY + (w.resultHeight-child.GetResultHeight())/2
			default:
				panic("not implemented")
			}

			child.UpdatePosition(currentX, currentY)

			currentX += child.GetResultWidth()
			currentX += w.op.Gap
		}
	} else if w.op.LayoutDirection == LayoutColumn {
		currentY := 0
		if w.op.AlignVertical == AlignStart {
			currentY = w.resultY + w.op.Padding.Top
		} else if w.op.AlignVertical == AlignCenter {
			childHeight := w.op.Gap * (len(w.Children) - 1)

			for _, child := range w.Children {
				if child.GetPositionMode() == PositionFixed || child.GetPositionMode() == PositionRelative {
					childHeight -= w.op.Gap
					continue
				}
				childHeight += child.GetResultHeight()
			}

			currentY = w.resultY + (w.resultHeight-childHeight)/2

		} else {
			panic("unimplemented")
		}

		for _, child := range w.Children {
			if child.GetPositionMode() == PositionFixed || child.GetPositionMode() == PositionRelative {
				child.UpdatePosition(w.resultX, w.resultY)
				continue
			}
			currentX := 0
			switch w.op.AlignHorizontal {
			case AlignStart:
				currentX = w.resultX + w.op.Padding.Left
			case AlignCenter:
				currentX = w.resultX + (w.GetResultWidth()-child.GetResultWidth())/2
			default:
				panic("unimplemented")
			}
			child.UpdatePosition(currentX, currentY)

			currentY += child.GetResultHeight()
			currentY += w.op.Gap
		}
	}
}

// generic things

func (w *Box) GetResultWidth() int {
	w.panicIfLayoutNotComputed()
	return w.resultWidth
}

func (w *Box) GetResultHeight() int {
	w.panicIfLayoutNotComputed()
	return w.resultHeight
}

func (w *Box) GetResultX() int {
	w.panicIfLayoutNotComputed()
	return w.resultX
}

func (w *Box) GetResultY() int {
	w.panicIfLayoutNotComputed()
	return w.resultY
}

func (w *Box) CanGrowWidth() bool {
	return w.op.WidthMode == SizeGrow
}
func (w *Box) CanGrowHeight() bool {
	return w.op.HeightMode == SizeGrow
}

func (w *Box) GetMinWidth() int {
	w.panicIfLayoutNotComputed()
	return w.minWidth
}
func (w *Box) GetMinHeight() int {
	w.panicIfLayoutNotComputed()
	return w.minHeight
}

func (w *Box) SetResultWidth(value int) {
	w.resultWidth = value
}
func (w *Box) SetResultHeight(value int) {
	w.resultHeight = value
}

func (w *Box) String() string {
	return w.string("")
}

func (w *Box) string(padding string) string {
	str := padding + fmt.Sprint(
		"<Box",
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
		case *Box:
			str += child.string(padding+"  ") + "\n"
		default:
			str += padding + "  " + child.String() + "\n"
		}
	}
	str += padding + "</Box>"
	return str
}

func (w *Box) UpdateInput(ctx *UIContext) {
	ctx.elementCache[w.id] = w

	if w.op.CursorShape != 0 && ctx.isWidgetHovered(w) {
		ctx.SetCursorShape(w.op.CursorShape)
	}

	if ctx.isWidgetHovered(w) {
		ctx.hoveredWidget = w.GetUniqueId()
		ctx.hoveredWidgets = append(ctx.hoveredWidgets, w.GetUniqueId())
	}

	if w.postLayoutFunction != nil {
		w.postLayoutFunction()
	}

	for _, child := range w.Children {
		child.UpdateInput(ctx)
	}

}

func (w *Box) GetPositionMode() PositionMode {
	return w.op.PositionMode
}

// TODO: add ui element to top of stack trace?
func (w *Box) panicIfLayoutNotComputed() {
	if !w.ctx.IsLayoutComputedOrComputing() {
		panic("Layout not computed! Run layout dependent code .....")
	}
}
