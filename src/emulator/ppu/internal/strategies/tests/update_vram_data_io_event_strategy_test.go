package tests

import (
	"nes-emu/src/emulator/ppu/internal/domain"
	"nes-emu/src/emulator/ppu/internal/strategies"
	"nes-emu/test/fixtures"
	"testing"
)

func TestUpdateVRAMDataIOEventStrategy_PanicWhenNotInitialized(t *testing.T) {
	bus := fixtures.NewMockBus()
	sut := strategies.NewUpdateVRAMDataIOEventStrategy(bus)
	state := domain.NewPPU()
	sut.Handle(state)
}

func TestUpdateVRAMDataIOEventStrategy_Handle_WritesDataToVideoMemory(t *testing.T) {
	bus := fixtures.NewMockBus()
	sut := strategies.NewUpdateVRAMDataIOEventStrategy(bus)
	state := domain.NewPPU()

	targetVRAMAddress := uint16(0x2150)
	state.AppendOnV(0x21)
	state.AppendOnV(0x50)

	// Set data value in PPU_DATA register (0x2007)
	expectedValue := uint8(0x7F)
	bus.WriteToMemory(domain.PPU_DATA, expectedValue)

	sut.Handle(state)

	actualValue := bus.ReadFromVideoMemory(targetVRAMAddress)
	if actualValue != expectedValue {
		t.Errorf("expected video memory at 0x%04X to be 0x%02X, got 0x%02X", targetVRAMAddress, expectedValue, actualValue)
	}
}

func TestUpdateVRAMDataIOEventStrategy_Handle_IncrementsVBy1_WhenIncrementBitClear(t *testing.T) {
	testCases := []struct {
		name       string
		controlVal uint8
	}{
		{name: "control is 0", controlVal: 0x00},
		{name: "control has other bits set but bit 2 is 0", controlVal: 0b11111011},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bus := fixtures.NewMockBus()
			bus.WriteToMemory(domain.PPU_CONTROL, tc.controlVal)

			sut := strategies.NewUpdateVRAMDataIOEventStrategy(bus)
			state := domain.NewPPU()
			state.AppendOnV(0x20)
			state.AppendOnV(0x00)

			sut.Handle(state)

			if state.GetV() != 0x2001 {
				t.Errorf("expected V to be incremented by 1 (0x2001), got 0x%04X", state.GetV())
			}
		})
	}
}

func TestUpdateVRAMDataIOEventStrategy_Handle_IncrementsVBy32_WhenIncrementBitSet(t *testing.T) {
	testCases := []struct {
		name       string
		controlVal uint8
	}{
		{name: "only bit 2 is set (0x04)", controlVal: 0b00000100},
		{name: "all bits set (0xFF)", controlVal: 0xFF},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bus := fixtures.NewMockBus()
			bus.WriteToMemory(domain.PPU_CONTROL, tc.controlVal)

			sut := strategies.NewUpdateVRAMDataIOEventStrategy(bus)
			state := domain.NewPPU()
			state.AppendOnV(0x20)
			state.AppendOnV(0x00)

			sut.Handle(state)

			if state.GetV() != 0x2020 {
				t.Errorf("expected V to be incremented by 32 (0x2020), got 0x%04X", state.GetV())
			}
		})
	}
}
