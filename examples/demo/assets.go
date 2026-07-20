package main

import (
	"embed"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

var hud_button_nine_slice = LoadImageFromPath("assets/ui/hud_button_nine_slice.png")
var hud_button_nine_slice_inverted = LoadImageFromPath("assets/ui/hud_button_nine_slice_inverted.png")
var button_nine_slice = LoadImageFromPath("assets/ui/button_nine_slice.png")
var button_nine_slice_inverted = LoadImageFromPath("assets/ui/button_nine_slice_inverted.png")
var button_nine_slice_disabled = LoadImageFromPath("assets/ui/button_nine_slice_disabled.png")
var box_nine_slice = LoadImageFromPath("assets/ui/box_nine_slice.png")

var crafting_divider = LoadImageFromPath("assets/ui/crafting_divider.png")

var hotbar_slot = LoadImageFromPath("assets/ui/hotbar_slot.png")
var hotbar_slot_unselected = LoadImageFromPath("assets/ui/hotbar_slot_unselected.png")

var recipe_slot = LoadImageFromPath("assets/ui/recipe_slot.png")
var recipe_slot_active = LoadImageFromPath("assets/ui/recipe_slot_active.png")

var button_rectangle_depth_line = LoadImageFromPath("assets/kenney/button_rectangle_depth_line.png")
var button_rectangle_line = LoadImageFromPath("assets/kenney/button_rectangle_line.png")
var button_rectangle_line_pressed = LoadImageFromPath("assets/kenney/button_rectangle_line_pressed.png")
var button_rectangle_line_hover = LoadImageFromPath("assets/kenney/button_rectangle_line_hover.png")

var radio_deselected = LoadImageFromPath("assets/kenney/check_round_grey.png")
var radio_selected = LoadImageFromPath("assets/kenney/check_round_color.png")
var radio_hovered = LoadImageFromPath("assets/kenney/check_round_grey_hovered.png")

var input_outline_rectangle = LoadImageFromPath("assets/kenney/input_outline_rectangle.png")
var background_rectangle_line = LoadImageFromPath("assets/kenney/background_rectangle_line.png")
var check_square_grey = LoadImageFromPath("assets/kenney/check_square_grey.png")
var check_square_grey_checkmark = LoadImageFromPath("assets/kenney/check_square_grey_checkmark.png")

//go:embed assets/**
var emb embed.FS

func LoadImageFromPath(path string) *ebiten.Image {
	file, err := emb.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	sheet := ebiten.NewImageFromImage(img)

	return sheet
}
