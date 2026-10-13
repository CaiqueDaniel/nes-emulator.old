package strategies

import (
	"nes-emu/src/emulator/ppu/internal/domain"
	shared "nes-emu/src/emulator/shared/application"
)

// Simulate read to $2002 latch
type ResetWAndVBlankFlagIOEventStrategy struct {
	initialized bool
	bus         shared.MNIBus
}

func NewResetWAndVBlankFlagIOEventStrategy(bus shared.MNIBus) *ResetWAndVBlankFlagIOEventStrategy {
	return &ResetWAndVBlankFlagIOEventStrategy{
		initialized: true,
		bus:         bus,
	}
}

func (s *ResetWAndVBlankFlagIOEventStrategy) Handle(state *domain.PPU) {
	s.checkIfInitilized()

	status := s.bus.ReadFromMemory(domain.PPU_STATUS)
	s.bus.WriteToMemory(domain.PPU_STATUS, status&0x7F)
	state.ClearW()
}

func (s *ResetWAndVBlankFlagIOEventStrategy) checkIfInitilized() {
	if !s.initialized {
		panic("UpdateVRAMDataIOEventStrategy not initialized")
	}
}
