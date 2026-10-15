package strategies

import (
	"nes-emu/src/emulator/ppu/internal/domain"
	shared "nes-emu/src/emulator/shared/application"
)

// Simulate write to $2007 latch
type UpdateVRAMDataIOEventStrategy struct {
	initialized bool
	bus         shared.MNIBus
}

func NewUpdateVRAMDataIOEventStrategy(bus shared.MNIBus) *UpdateVRAMDataIOEventStrategy {
	return &UpdateVRAMDataIOEventStrategy{
		initialized: true,
		bus:         bus,
	}
}

func (s *UpdateVRAMDataIOEventStrategy) Handle(state *domain.PPU) {
	s.checkIfInitilized()

	value := s.bus.ReadFromMemory(domain.PPU_DATA)
	s.writeToVMemory(state.GetV(), value)
	s.increaseVByOffset(state)
}

func (p *UpdateVRAMDataIOEventStrategy) increaseVByOffset(state *domain.PPU) {
	if p.isOffsetIncrementBy32() {
		state.IncreaseVByOffset(32)
	} else {
		state.IncreaseVByOffset(1)
	}
}

func (p *UpdateVRAMDataIOEventStrategy) readVMemory(address uint16) uint8 {
	return p.bus.ReadFromVideoMemory(address)
}

func (p *UpdateVRAMDataIOEventStrategy) writeToVMemory(address uint16, value byte) {
	p.bus.WriteToVideoMemory(address, value)
}

func (p *UpdateVRAMDataIOEventStrategy) isOffsetIncrementBy32() bool {
	return p.bus.ReadFromMemory(domain.PPU_CONTROL)&0b100 != 0
}

func (s *UpdateVRAMDataIOEventStrategy) checkIfInitilized() {
	if !s.initialized {
		panic("UpdateVRAMDataIOEventStrategy not initialized")
	}
}
