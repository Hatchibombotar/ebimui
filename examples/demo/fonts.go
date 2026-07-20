package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	fontFaceRegular  = "assets/fonts/Lexend-Regular.ttf"
	fontFaceSemiBold = "assets/fonts/Lexend-SemiBold.ttf"
	fontFaceBold     = "assets/fonts/Lexend-Bold.ttf"
)

// TODO: These were *text.Face, now *text.GoTextFace. investigate
type fonts struct {
	face         *text.GoTextFace
	titleFace    *text.GoTextFace
	buttonFace   *text.GoTextFace
	tabFace      *text.GoTextFace
	bigTitleFace *text.GoTextFace
	toolTipFace  *text.GoTextFace
}

func loadFonts() (*fonts, error) {
	fontFace, err := loadFont(fontFaceRegular, 20)
	if err != nil {
		return nil, err
	}

	tabFace, err := loadFont(fontFaceSemiBold, 20)
	if err != nil {
		return nil, err
	}

	titleFontFace, err := loadFont(fontFaceBold, 24)
	if err != nil {
		return nil, err
	}

	bigTitleFontFace, err := loadFont(fontFaceBold, 28)
	if err != nil {
		return nil, err
	}

	toolTipFace, err := loadFont(fontFaceRegular, 15)
	if err != nil {
		return nil, err
	}
	buttonFace, err := loadFont(fontFaceSemiBold, 20)
	if err != nil {
		return nil, err
	}

	return &fonts{
		face:         fontFace,
		titleFace:    titleFontFace,
		bigTitleFace: bigTitleFontFace,
		toolTipFace:  toolTipFace,
		tabFace:      tabFace,
		buttonFace:   buttonFace,
	}, nil
}

func loadFont(path string, size float64) (*text.GoTextFace, error) {
	fontFile, err := emb.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := text.NewGoTextFaceSource(fontFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}

var loadedFonts *fonts

func init() {
	fonts, err := loadFonts()
	if err != nil {
		panic(err)
	}
	loadedFonts = fonts
}
