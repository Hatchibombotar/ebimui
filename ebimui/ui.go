package ebimui

import "github.com/hajimehoshi/ebiten/v2"

type GenericWidget interface {
	Draw(screen *ebiten.Image, root *UIContext)
	Update()

	// TODO: Consider removing. Change cursor globally using hover widget.
	UpdateInput(*UIContext)

	// Prefix all non init stuff with Layout___ ?
	UpdateSizeWidthFitPass()
	UpdateSizeHeightFitPass()
	UpdateSizeWrapWidth()
	UpdateSizeWidthGrowPass()
	UpdateSizeHeightGrowPass()
	UpdatePosition(x, y int)

	// SetParent(parent GenericWidget)
	GetResultWidth() int
	GetResultHeight() int
	GetResultX() int
	GetResultY() int
	CanGrowWidth() bool
	CanGrowHeight() bool
	GetMinWidth() int
	GetMinHeight() int

	SetResultWidth(int)
	SetResultHeight(int)

	String() string

	// TODO: Consider if PositionMode should be seperate from "floating" e.g. not affecting parent widgets
	// This may allow us to remove GetPositionMode function.
	GetPositionMode() PositionMode

	GetUniqueId() WidgetID
}

type WidgetID struct {
	callerId       uintptr
	callerInstance int
}

type SizeMode = int

const (
	SizeFit SizeMode = iota
	SizeFixed
	SizeGrow
)

type LayoutDirection int

const (
	LayoutColumn LayoutDirection = iota
	LayoutRow
)

type LayoutType int

const (
	LayoutFlex LayoutType = iota
	LayoutGrid
)

type AlignMode int

const (
	AlignStart AlignMode = iota
	AlignCenter
	AlignEnd
)

type PositionMode = int

const (
	PositionAuto PositionMode = iota
	PositionFixed
	PositionRelative
)

type Spacing struct {
	Top, Right, Bottom, Left int
}

func (s Spacing) WithAll(value int) Spacing {
	s.Top = value
	s.Right = value
	s.Bottom = value
	s.Left = value

	return s
}
