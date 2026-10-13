package services

import (
	"image"
	"image/color"
	"nes-emu/src/emulator/ppu/internal/domain"

	"golang.org/x/exp/shiny/screen"
)

type ShinyScreen struct {
	window *screen.Window
	buffer *screen.Buffer
}

func NewShinyScreen(window *screen.Window, buffer *screen.Buffer) *ShinyScreen {
	return &ShinyScreen{
		window: window,
		buffer: buffer,
	}
}

func (s *ShinyScreen) ShowImage(buffer *[domain.MAX_FRAME_SCANLINE + 1][domain.MAX_PIXEL_PER_SCANLINE + 1]uint32) {
	img := (*s.buffer).RGBA()

	for y, scanline := range *buffer {
		for x, pixel := range scanline {
			c := color.RGBA{
				R: uint8(pixel >> 16),
				G: uint8(pixel >> 8),
				B: uint8(pixel),
				A: 255,
			}
			img.SetRGBA(x, y, c)
		}
	}

	(*s.window).Upload(image.Point{0, 0}, *s.buffer, (*s.buffer).Bounds())
	(*s.window).Publish()
}
