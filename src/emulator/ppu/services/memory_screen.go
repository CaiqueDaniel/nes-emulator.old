package services

import (
	"nes-emu/src/emulator/ppu/internal/domain"
)

type MemoryScreen struct {
	buffer *[domain.MAX_FRAME_SCANLINE + 1][domain.MAX_PIXEL_PER_SCANLINE + 1]uint32
}

func NewMemoryScreen() *MemoryScreen {
	return &MemoryScreen{}
}

func (s *MemoryScreen) ShowImage(buffer *[domain.MAX_FRAME_SCANLINE + 1][domain.MAX_PIXEL_PER_SCANLINE + 1]uint32) {
	s.buffer = buffer
}

func (s *MemoryScreen) GetImageBuffer() *[domain.MAX_FRAME_SCANLINE + 1][domain.MAX_PIXEL_PER_SCANLINE + 1]uint32 {
	return s.buffer
}
