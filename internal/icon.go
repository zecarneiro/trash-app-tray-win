package internal

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/energye/systray"
)

// Matriz Pixel-Art exactly of the icon (16x16)
var baseGrid = [16][16]int{
	{0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
	{0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0},
	{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	{0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0},
	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
}

// pngToIco packages a PNG image in Windows icon format (.ico) in memory.
func pngToIco(pngBytes []byte) []byte {
	var icoBuf bytes.Buffer
	// ICONDIR Header (6 bytes)
	icoBuf.Write([]byte{
		0, 0, // Reserved (deve ser 0)
		1, 0, // Type: 1 = .ICO
		1, 0, // Images number: 1
	})
	size := len(pngBytes)
	// ICONDIRENTRY (16 bytes)
	icoBuf.Write([]byte{
		32,   // Width (32px)
		32,   // Height (32px)
		0,    // Color Count (0 para PNG/24bit+)
		0,    // Reserved
		1, 0, // Color planes (1)
		32, 0, // Bits by pixel (32 bpp)
		byte(size), byte(size >> 8), byte(size >> 16), byte(size >> 24), // Data size in bytes
		22, 0, 0, 0, // Offset where the PNG data begins (22 = 6 header + 16 entry)
	})
	// Append the PNG bytes after the header.
	icoBuf.Write(pngBytes)
	return icoBuf.Bytes()
}

// Generate image in bytes (PNG) with escale Pixel-Perfect 32x32
func generateIconData(isFull bool) []byte {
	grid := baseGrid
	if isFull {
		// Interior filling
		for r := 5; r <= 13; r++ {
			for c := 4; c <= 11; c++ {
				grid[r][c] = 2
			}
		}
		// Wave pattern on top of the liquid
		grid[4][4] = 2
		grid[4][5] = 0
		grid[4][6] = 2
		grid[4][7] = 2
		grid[4][8] = 0
		grid[4][9] = 2
		grid[4][10] = 2
		grid[4][11] = 0
	}
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))

	// Colors
	white := color.RGBA{255, 255, 255, 255}
	cyan := color.RGBA{0, 230, 255, 255}
	transparent := color.RGBA{0, 0, 0, 0}

	for y := 0; y < 16; y++ {
		for x := 0; x < 16; x++ {
			var pxColor color.Color
			switch grid[y][x] {
			case 1:
				pxColor = white
			case 2:
				pxColor = cyan
			default:
				pxColor = transparent
			}
			// Renders each pixel in a 2x2 block for perfect sharpness.
			for dy := 0; dy < 2; dy++ {
				for dx := 0; dx < 2; dx++ {
					img.Set(x*2+dx, y*2+dy, pxColor)
				}
			}
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return pngToIco(buf.Bytes())
}

func buildIconMenu() {
	saveIncon := systray.AddMenuItem("Save Icon", "Save Icon")
	saveIncon.Click(func() {
		err := os.WriteFile("favicon.ico", generateIconData(false), 0644)
		if err != nil {
			fmt.Println("Erro ao guardar o ícone:", err)
			return
		}
	})
	saveIncon.Disable()
	saveIncon.Hide()
}
