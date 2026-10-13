package tests

import (
	"nes-emu/src/emulator/ppu/internal/domain"
	"nes-emu/src/emulator/ppu/internal/strategies"
	"nes-emu/test/fixtures"
	"testing"
)

func TestUpdateVRAMAddressIOEventStrategy_PanicWhenNotInitialized(t *testing.T) {
	bus := fixtures.NewMockBus()
	sut := strategies.NewUpdateVRAMAddressIOEventStrategy(bus)
	state := domain.NewPPU()
	sut.Handle(state)
}

func TestUpdateVRAMAddressIOEventStrategy_Handle_FirstWrite(t *testing.T) {
	bus := fixtures.NewMockBus()
	sut := strategies.NewUpdateVRAMAddressIOEventStrategy(bus)
	state := domain.NewPPU()

	// Write 0x3D to PPU_ADDRESS
	bus.WriteToMemory(domain.PPU_ADDRESS, 0x3D)

	sut.Handle(state)

	// (0x3D << 8) & 0x3F00 = 0x3D00
	if state.GetT() != 0x3D00 {
		t.Errorf("expected T to be 0x3D00, got 0x%04X", state.GetT())
	}

	if !state.IsWSet() {
		t.Errorf("expected W to be toggled to true after first write")
	}

	if state.GetV() != 0 {
		t.Errorf("expected V to remain 0 after first write, got 0x%04X", state.GetV())
	}
}

func TestUpdateVRAMAddressIOEventStrategy_Handle_FirstWrite_MasksBitsAbove14(t *testing.T) {
	bus := fixtures.NewMockBus()
	sut := strategies.NewUpdateVRAMAddressIOEventStrategy(bus)
	state := domain.NewPPU()

	// Write 0xFF to PPU_ADDRESS: bits 6 and 7 should be cleared by 0x3F00 mask
	bus.WriteToMemory(domain.PPU_ADDRESS, 0xFF)

	sut.Handle(state)

	// (0xFF << 8) & 0x3F00 = 0x3F00
	if state.GetT() != 0x3F00 {
		t.Errorf("expected T to be masked to 0x3F00, got 0x%04X", state.GetT())
	}
}

func TestUpdateVRAMAddressIOEventStrategy_Handle_SecondWrite(t *testing.T) {
	bus := fixtures.NewMockBus()
	sut := strategies.NewUpdateVRAMAddressIOEventStrategy(bus)
	state := domain.NewPPU()

	// Prepare state as if first write already occurred
	state.SetT(0x2100)
	state.ToggleW() // W is now true
	if !state.IsWSet() {
		t.Fatal("expected W to be true")
	}

	// Write low byte to PPU_ADDRESS
	bus.WriteToMemory(domain.PPU_ADDRESS, 0x08)

	sut.Handle(state)

	expectedAddress := uint16(0x2108)
	if state.GetT() != expectedAddress {
		t.Errorf("expected T to be 0x%04X, got 0x%04X", expectedAddress, state.GetT())
	}

	if state.GetV() != expectedAddress {
		t.Errorf("expected V to be copied from T (0x%04X), got 0x%04X", expectedAddress, state.GetV())
	}

	if state.IsWSet() {
		t.Errorf("expected W to be toggled back to false after second write")
	}
}

func TestUpdateVRAMAddressIOEventStrategy_Handle_FullTwoWriteCycle(t *testing.T) {
	bus := fixtures.NewMockBus()
	sut := strategies.NewUpdateVRAMAddressIOEventStrategy(bus)
	state := domain.NewPPU()

	// 1. First write: high byte 0x24
	bus.WriteToMemory(domain.PPU_ADDRESS, 0x24)
	sut.Handle(state)

	if state.GetT() != 0x2400 {
		t.Fatalf("expected T to be 0x2400 after 1st write, got 0x%04X", state.GetT())
	}
	if !state.IsWSet() {
		t.Fatalf("expected W to be true after 1st write")
	}
	if state.GetV() != 0x0000 {
		t.Fatalf("expected V to be 0 after 1st write, got 0x%04X", state.GetV())
	}

	// 2. Second write: low byte 0x12
	bus.WriteToMemory(domain.PPU_ADDRESS, 0x12)
	sut.Handle(state)

	if state.GetT() != 0x2412 {
		t.Fatalf("expected T to be 0x2412 after 2nd write, got 0x%04X", state.GetT())
	}
	if state.GetV() != 0x2412 {
		t.Fatalf("expected V to be 0x2412 after 2nd write, got 0x%04X", state.GetV())
	}
	if state.IsWSet() {
		t.Fatalf("expected W to be false after 2nd write")
	}
}
