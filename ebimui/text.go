package ebimui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Text struct {
	text string
	op   *TextWidgetOps
	// parent GenericWidget

	resultX, resultY          int
	resultWidth, resultHeight int
	minWidth, minHeight       int

	// resultFontFace *text.GoTextFace

	lines []string

	id WidgetID

	ctx *UIContext
}

func (t *Text) GetUniqueId() WidgetID {
	return t.id
}

type WrapBehaviour int

const (
	NoWrap WrapBehaviour = iota
	WrapText
)

type TextWidgetOps struct {
	Color         color.Color
	Face          text.Face
	WrapBehaviour WrapBehaviour

	TextAlign AlignMode
}

func (w *Text) Draw(screen *ebiten.Image, root *UIContext) {
	face := w.op.Face

	for i, line := range w.lines {
		op := &text.DrawOptions{}

		fontSize := face.Metrics().HAscent

		if w.op.TextAlign == AlignCenter {
			lineWidth := textWidth(face, line)
			op.GeoM.Translate(float64(w.resultWidth)/2-float64(lineWidth)/2, 0)
		}

		op.GeoM.Translate(float64(w.resultX), float64(w.resultY+(i*int(fontSize))))

		if w.op.Color == nil {
			panic("Colour not defined for text")
		}
		op.ColorScale.ScaleWithColor(w.op.Color)

		text.Draw(screen, line, face, op)
	}
}

func (w *Text) Update() {
	panic("Text Update() function should never be called.")
}
func (w *Text) UpdateSizeWidthFitPass() {
	minWidth, maxWidth := w.calculateTextBounds()
	w.minWidth = minWidth
	w.resultWidth = maxWidth

	w.minHeight = 0
	w.resultHeight = 0
}
func (w *Text) UpdateSizeHeightFitPass()  {}
func (w *Text) UpdateSizeWidthGrowPass()  {}
func (w *Text) UpdateSizeHeightGrowPass() {}

func (w *Text) UpdateSizeWrapWidth() {
	w.calculateLineWrappings()
}

func (w *Text) UpdatePosition(x, y int) {
	w.resultX = x
	w.resultY = y
}

func (w *Text) calculateTextBounds() (minWidth int, maxWidth int) {
	if w.op.Face == nil {
		panic("font face not definied")
	}
	if w.op.WrapBehaviour == NoWrap {
		width := textWidth(w.op.Face, w.text)
		return width, width
	}
	maxTextWidth := 0
	totalWidth := 0

	startIndex := 0
	for endIndex, char := range w.text {
		if char == ' ' || char == '\n' || char == '\t' {
			str := string(w.text[startIndex:endIndex])
			wordWidth := textWidth(w.op.Face, str)

			totalWidth += wordWidth
			if char == ' ' || char == '\t' {
				totalWidth += textWidth(w.op.Face, string(char))
			}

			maxTextWidth = max(maxTextWidth, wordWidth)
			startIndex = endIndex
		}
	}
	str := string(w.text[startIndex:len(w.text)])
	wordWidth := textWidth(w.op.Face, str)

	totalWidth += wordWidth

	maxTextWidth = max(maxTextWidth, wordWidth)

	return maxTextWidth, totalWidth
}

func textWidth(fontFace text.Face, str string) int {
	w, _ := text.Measure(str, fontFace, 0)
	return int(w)
	// return int(text.Advance(str, fontFace))
}

func (w *Text) calculateLineWrappings() {
	face := w.op.Face
	lines := make([]string, 0)

	if w.op.WrapBehaviour == NoWrap {
		w.resultWidth = textWidth(face, w.text)
		// TODO: make all spacing adhere to metrics
		fontSize := face.Metrics().HAscent + face.Metrics().HDescent
		w.resultHeight = int(fontSize)
		lines = append(lines, w.text)
		w.lines = lines
		return
	}
	var endIdx, p int
	maxWidth := 0

	for endIdx < len(w.text) {
		wi := 0
		endIdx = p
		startIdx := endIdx
		for endIdx < len(w.text) && w.text[endIdx] != '\n' {
			word := p
			for p < len(w.text) && w.text[p] != ' ' && w.text[p] != '\n' {
				p++
			}
			if wi > maxWidth {
				maxWidth = wi
			}
			wi += textWidth(face, w.text[word:p])
			if wi > w.resultWidth && endIdx != startIdx {
				break
			}
			if p < len(w.text) {
				wi += textWidth(face, string(w.text[p]))
			}
			endIdx = p
			p++
		}

		lines = append(lines, w.text[startIdx:endIdx])
		p = endIdx + 1
	}

	w.lines = lines
	// if maxWidth < w.resultWidth {
	// 	w.InternalWidget.width = maxWidth
	// } else {
	// 	w.InternalWidget.width = w.resultWidth
	// }
	fontSize := face.Metrics().HAscent
	w.resultHeight = len(lines) * int(fontSize)
	if len(lines) > 1 {
		// previously removed for some reason
		// reason: when added it breaks shit
		w.resultWidth = maxWidth
	}
}

func (w *Text) GetResultWidth() int {
	return w.resultWidth
}

func (w *Text) GetResultHeight() int {
	return w.resultHeight
}

func (w *Text) GetResultX() int {
	return w.resultX
}

func (w *Text) GetResultY() int {
	return w.resultY
}

func (w *Text) CanGrowWidth() bool {
	return false
}
func (w *Text) CanGrowHeight() bool {
	return false
}

func (w *Text) GetMinHeight() int {
	return 0
}
func (w *Text) GetMinWidth() int {
	// if w.Op.WrapBehaviour == NoWrap {
	width, _ := w.calculateTextBounds()
	return width
	// }
	// return 0
}

func (w *Text) SetResultWidth(value int) {
	w.resultWidth = value
}
func (w *Text) SetResultHeight(value int) {
	w.resultHeight = value
}

func (w *Text) String() string {
	return fmt.Sprint(
		"<Text",
		" id=",
		w.GetUniqueId(),
		" width=",
		w.resultWidth,
		" height=",
		w.resultHeight,
		" lines=",
		len(w.lines),
		" size=",
		w.op.Face.Metrics().HAscent,
		"/>",
	)
}

func (w *Text) UpdateInput(ctx *UIContext) {
	ctx.elementCache[w.id] = w

	if ctx.isWidgetHovered(w) {
		ctx.hoveredWidget = w.GetUniqueId()
	}
}

func (w *Text) SetText(newText string) {
	w.text = newText
}

func (w *Text) GetPositionMode() PositionMode {
	return PositionAuto
}
