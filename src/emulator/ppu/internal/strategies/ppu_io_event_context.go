package strategies

import (
	"nes-emu/src/emulator/ppu/internal/application"
	"nes-emu/src/emulator/ppu/internal/domain"
	shared "nes-emu/src/emulator/shared/application"
)

type PPUIOEventContext struct {
	initialized bool
	bus         shared.MNIBus
}

func NewPPUIOEventContext(bus shared.MNIBus) *PPUIOEventContext {
	return &PPUIOEventContext{
		initialized: true,
		bus:         bus,
	}
}

func (p *PPUIOEventContext) HandleEvent(event *application.IOEvent, state *domain.PPU) {
	if !p.initialized {
		panic("PPUIOEventContext not initialized")
	}

	var context application.IOEventsStrategy

	if event.IsWrite {
		context = p.getStrategyForWriteSignal(event.Address, state)
	} else {
		context = p.getStrategyForReadSignal(event.Address, state)
	}

	if context != nil {
		context.Handle(state)
	}
}

func (p *PPUIOEventContext) getStrategyForWriteSignal(address uint16, state *domain.PPU) application.IOEventsStrategy {
	switch address {
	case domain.PPU_ADDRESS:
		return NewUpdateVRAMAddressIOEventStrategy(p.bus)

	case domain.PPU_DATA:
		return NewUpdateVRAMDataIOEventStrategy(p.bus)

	case domain.PPU_SCROLL:
		return NewCopyScrollToVRegister(p.bus)
	}

	return nil
}

func (p *PPUIOEventContext) getStrategyForReadSignal(address uint16, state *domain.PPU) application.IOEventsStrategy {
	switch address {
	case domain.PPU_STATUS:
		return NewResetWAndVBlankFlagIOEventStrategy(p.bus)
	}

	return nil
}
