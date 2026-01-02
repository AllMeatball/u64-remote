package screen

import (
	// "fmt"
	"bytes"
	_ "embed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"github.com/AllMeatball/u64-remote/server"
)

//go:embed petscii-upper.png
var PETSCII_UPPER_PNG []byte

var PETSCII_UPPER *image.Paletted

const ROWS = 25
const COLS = 40
const CHAR_SIZE = 8

const (
	BLACK = iota
	WHITE
	RED
	CYAN
	PURPLE
	GREEN
	BLUE
	YELLOW
	ORANGE
	BROWN
	LIGHT_RED
	DARK_GRAY
	GRAY
	LIGHT_GREEN
	LIGHT_BLUE
	LIGHT_GRAY
)

var COLORS = [16]color.RGBA {
	{0x00, 0x00, 0x00, 0xFF},
	{0xFF, 0xFF, 0xFF, 0xFF},
	{0x68, 0x37, 0x2B, 0xFF},
	{0x70, 0xA4, 0xB2, 0xFF},
	{0x6F, 0x3D, 0x86, 0xFF},
	{0x58, 0x8D, 0x43, 0xFF},
	{0x35, 0x28, 0x79, 0xFF},
	{0xB8, 0xC7, 0x6F, 0xFF},
	{0x6F, 0x4F, 0x25, 0xFF},
	{0x43, 0x39, 0x00, 0xFF},
	{0x9A, 0x67, 0x59, 0xFF},
	{0x44, 0x44, 0x44, 0xFF},
	{0x6C, 0x6C, 0x6C, 0xFF},
	{0x9A, 0xD2, 0x84, 0xFF},
	{0x6C, 0x5E, 0xB5, 0xFF},
	{0x95, 0x95, 0x95, 0xFF},
}

var IMAGE_PALETTE = []color.Color{}

func Setup() error {
	var err error

	for _, color := range COLORS {
		IMAGE_PALETTE = append(IMAGE_PALETTE, color)
	}

	img, err := png.Decode(bytes.NewBuffer(PETSCII_UPPER_PNG))
	if err != nil { return err }

	PETSCII_UPPER = img.(*image.Paletted)

	return nil
}

func GetCharImage(char, fg, bg byte) image.Image {
	char = 1
	var x int = (int(char) % 16) * CHAR_SIZE
	var y int = (int(char) / 16) * CHAR_SIZE

	rect := image.Rect(x, y, x+CHAR_SIZE, y+CHAR_SIZE)

	// fmt.Printf("%02x\n", color)
	fg_color := COLORS[1]
	bg_color := COLORS[0]

	char_img := PETSCII_UPPER.SubImage(rect)
	char_paletted := char_img.(*image.Paletted)

	char_paletted.Palette[0] = bg_color
	char_paletted.Palette[1] = fg_color

	// image
	return char_img
}

func GrabPetsciiScreen(server server.U64Server) image.Image {
	// chars, err := server.PeekMemory(0x0400, COLS)
	// if err != nil { panic(err) }

	colors, err := server.PeekMemory(0xD800, COLS)
	if err != nil { panic(err) }

	bg_bytes, err := server.PeekMemory(0xD021, 1)
	if err != nil { panic(err) }

	bg_color := bg_bytes[0]
	// fmt.Printf("%02x\n", bg_color)

	img := image.NewRGBA(image.Rect(0, 0, COLS * CHAR_SIZE, ROWS * CHAR_SIZE))


	for x := range COLS {
		// offset := image.Pt(x * 8, 0)
		// fmt.Println(chars[x])
		char_img := GetCharImage(0, colors[x], bg_color)
		draw.Draw(img, char_img.Bounds(), char_img, image.Point{}, draw.Over)
	}

	return img
}


