package strategies

import (
	"nes-emu/src/emulator/ppu/internal/domain"
	shared "nes-emu/src/emulator/shared/application"
)

type CopyScrollToVRegister struct {
	initialized bool
	bus         shared.MNIBus
}

func NewCopyScrollToVRegister(bus shared.MNIBus) *CopyScrollToVRegister {
	return &CopyScrollToVRegister{
		initialized: true,
		bus:         bus,
	}
}

func (s *CopyScrollToVRegister) Handle(state *domain.PPU) {
	s.checkInitialized()
	state.AppendOnV(s.bus.ReadFromMemory(domain.PPU_SCROLL))
}

func (s *CopyScrollToVRegister) checkInitialized() {
	if !s.initialized {
		panic("CopyScrollXAndYToScrollRegister strategy is not initialized")
	}
}
