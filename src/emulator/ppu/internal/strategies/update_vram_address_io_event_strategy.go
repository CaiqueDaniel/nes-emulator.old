package strategies

import (
	"nes-emu/src/emulator/ppu/internal/domain"
	shared "nes-emu/src/emulator/shared/application"
)

// Simulate write to $2006 latch
type UpdateVRAMAddressIOEventStrategy struct {
	initilized bool
	bus        shared.MNIBus
}

func NewUpdateVRAMAddressIOEventStrategy(bus shared.MNIBus) *UpdateVRAMAddressIOEventStrategy {
	return &UpdateVRAMAddressIOEventStrategy{
		initilized: true,
		bus:        bus,
	}
}

func (s *UpdateVRAMAddressIOEventStrategy) Handle(state *domain.PPU) {
	s.checkIfInitilized()
	state.AppendOnV(s.bus.ReadFromMemory(domain.PPU_ADDRESS))
}

func (s *UpdateVRAMAddressIOEventStrategy) checkIfInitilized() {
	if !s.initilized {
		panic("UpdateVRAMAddressIOEventStrategy not initialized")
	}
}
