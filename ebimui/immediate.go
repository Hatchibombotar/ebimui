package ebimui

import (
	"image/color"
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func caller() uintptr {
	var pcs [1]uintptr
	n := runtime.Callers(3, pcs[:])
	if n == 0 {
		return 0
	}
	return pcs[0]
}

func (uiContext *UIContext) NewBoxWidget(options func(b *Box)) *Box {
	box := &Box{
		op:  &BoxWidgetOps{},
		id:  uiContext.GetId(),
		ctx: uiContext,
	}

	options(box)

	return box
}

// Rename to AddNewBoxWidget?
func (box *Box) AppendNewBoxWidget(options func(b *Box)) {
	ctx := box.ctx
	box.AddChild(
		ctx.NewBoxWidget(options),
	)
}
func (box *Box) AppendNewTextWidget(options func(t *Text)) {
	ctx := box.ctx
	box.AddChild(
		ctx.NewTextWidget(options),
	)
}
func (box *Box) LayoutDirection(direction LayoutDirection) {
	box.op.LayoutDirection = direction
}

// TODO: rework how layout + columns work as columns only make sense for grid
func (box *Box) LayoutMode(mode LayoutType) {
	box.op.LayoutType = mode
}
func (box *Box) Columns(columns int) {
	box.op.Columns = columns
}

func (box *Box) Padding(Top, Right, Bottom, Left int) {
	box.op.Padding = Spacing{Top, Right, Bottom, Left}
}

func (box *Box) FixedWidth(Width int) {
	box.op.WidthMode = SizeFixed
	box.op.Width = Width
}

func (box *Box) FixedHeight(Height int) {
	box.op.HeightMode = SizeFixed
	box.op.Height = Height
}

func (box *Box) WidthGrow() {
	box.op.WidthMode = SizeGrow
}

func (box *Box) Gap(Gap int) {
	box.op.Gap = Gap
}

func (box *Box) AlignHorizontal(align AlignMode) {
	box.op.AlignHorizontal = align
}
func (box *Box) AlignVertical(align AlignMode) {
	box.op.AlignVertical = align
}
func (box *Box) PositionFixed(x, y int) {
	box.op.PositionMode = PositionFixed
	box.op.X = x
	box.op.Y = y
}
func (box *Box) PositionRelative(x, y int) {
	box.op.PositionMode = PositionRelative
	box.op.X = x
	box.op.Y = y
}

// Consider removing root
func (box *Box) OnDraw(drawFunc func(screen *ebiten.Image, widget GenericWidget)) {
	box.op.OnDraw = drawFunc
}
func (box *Box) DrawNineSlice(nineSliceImage *ebiten.Image, nineSliceWidth int) {
	box.op.OnDraw = func(screen *ebiten.Image, widget GenericWidget) {
		FillNineSlice(screen, widget, nineSliceImage, nineSliceWidth, nineSliceWidth, nineSliceWidth, nineSliceWidth)
	}
}
func (box *Box) DrawNineSlice4(nineSliceImage *ebiten.Image, Top, Right, Bottom, Left int) {
	box.op.OnDraw = func(screen *ebiten.Image, widget GenericWidget) {
		FillNineSlice(screen, widget, nineSliceImage, Top, Right, Bottom, Left)
	}
}

func (box *Box) DrawFillSolid(color color.Color) {
	box.op.OnDraw = func(screen *ebiten.Image, widget GenericWidget) {
		FillSolid(screen, widget, color)
	}
}

func (box *Box) CursorShape(cursorShape ebiten.CursorShapeType) {
	box.op.CursorShape = cursorShape
}

func (box *Box) DeferToPostLayout(f func()) {
	if box.postLayoutFunction != nil {
		panic("DeferToPostLayout already called")
	}
	box.postLayoutFunction = f
}

func (uiContext *UIContext) NewTextWidget(options func(t *Text)) *Text {
	text := &Text{
		op:  &TextWidgetOps{},
		id:  uiContext.GetId(),
		ctx: uiContext,
	}

	options(text)

	return text
}

func (text *Text) Content(content string) {
	text.text = content
}

func (text *Text) Wrap(wrapMode WrapBehaviour) {
	text.op.WrapBehaviour = wrapMode
}
func (text *Text) Face(face *text.GoTextFace) {
	text.op.Face = face
}

func (text *Text) Color(color color.Color) {
	text.op.Color = color
}
