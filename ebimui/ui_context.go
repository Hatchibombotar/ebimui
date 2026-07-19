package ebimui

import (
	"fmt"
	"image"
	"image/color"
	"runtime"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type UIContext struct {
	GetCursorPositionFunc func() (int, int)
	cursorShape           ebiten.CursorShapeType
	isHovered             bool
	currentUis            []GenericWidget

	elementInstanceCount map[uintptr]int
	elementCache         map[WidgetID]GenericWidget

	hoveredWidget  WidgetID
	hoveredWidgets []WidgetID

	hotWidget WidgetID

	layoutComputed bool
}

func (r *UIContext) GetId() WidgetID {
	var pcs [1]uintptr
	n := runtime.Callers(3, pcs[:])
	if n == 0 {
		fmt.Println("Failed to give id to widget.")
		return WidgetID{
			callerId:       0,
			callerInstance: 0,
		}
	}
	caller := pcs[0]

	count := r.elementInstanceCount[caller]

	r.elementInstanceCount[caller] += 1

	return WidgetID{
		callerId:       caller,
		callerInstance: count,
	}
}

func (r *UIContext) AddUi(widget GenericWidget) {
	r.currentUis = append(r.currentUis, widget)
}

func (r *UIContext) Update() {
	r.hoveredWidgets = []WidgetID{}
	r.hoveredWidget = WidgetID{}

	r.layoutComputed = true

	for _, ui := range r.currentUis {
		ui.Update()
	}
	for _, ui := range r.currentUis {
		ui.UpdateInput(r)
	}

	if r.hoveredWidget.callerId != 0 {
		r.isHovered = true
	}

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		r.hotWidget = WidgetID{}
	}

	ebiten.SetCursorShape(r.cursorShape)
}
func (ctx *UIContext) DrawChildren(screen *ebiten.Image) {
	// fmt.Println("draw", len(r.currentUis))
	for _, ui := range ctx.currentUis {
		ui.Draw(screen, ctx)
	}
	// ebitenutil.DebugPrint(screen, fmt.Sprintln(ebiten.ActualFPS(), ebiten.ActualTPS()))
	if ebiten.IsKeyPressed(ebiten.KeyF12) {
		debugWidgetId := ctx.hoveredWidget

		debugWidget, exists := ctx.elementCache[debugWidgetId]

		if exists {
			vector.FillRect(
				screen,
				float32(debugWidget.GetResultX()),
				float32(debugWidget.GetResultY()),
				float32(debugWidget.GetResultWidth()),
				float32(debugWidget.GetResultHeight()),
				color.RGBA{50, 50, 255, 150},
				false,
			)
			ebitenutil.DebugPrint(screen, debugWidget.String())
		}
	}
}

func (r *UIContext) IsHovered() bool {
	return r.isHovered
}

func (r *UIContext) IsWidgetHovered(widget GenericWidget) bool {
	id := widget.GetUniqueId()
	return slices.Contains(r.hoveredWidgets, id)
}

func (e *UIContext) isWidgetHovered(widget GenericWidget) bool {
	mouseX, mouseY := e.CursorPosition()
	r := image.Rect(widget.GetResultX(), widget.GetResultY(), widget.GetResultX()+widget.GetResultWidth(), widget.GetResultY()+widget.GetResultHeight())
	p := image.Point{mouseX, mouseY}

	return p.In(r)
}

func (r *UIContext) SetCursorShape(shape ebiten.CursorShapeType) {
	r.cursorShape = shape
}

func NewUIContext() *UIContext {
	root := &UIContext{
		elementCache:         map[WidgetID]GenericWidget{},
		elementInstanceCount: map[uintptr]int{},
	}

	return root
}

func (r *UIContext) CursorPosition() (int, int) {
	if r.GetCursorPositionFunc == nil {
		return ebiten.CursorPosition()
	} else {
		return r.GetCursorPositionFunc()
	}
}

func (r *UIContext) PreUpdate() {
	r.cursorShape = ebiten.CursorShapeDefault
	r.isHovered = false

	r.currentUis = []GenericWidget{}
	r.elementInstanceCount = map[uintptr]int{}

	r.layoutComputed = false
}

func (r *UIContext) IsLayoutComputedOrComputing() bool {
	return r.layoutComputed
}

func (r *UIContext) SetWidgetHot(widget GenericWidget) {
	r.hotWidget = widget.GetUniqueId()
}

func (r *UIContext) IsWidgetHot(widget GenericWidget) bool {
	return r.hotWidget == widget.GetUniqueId()
}
