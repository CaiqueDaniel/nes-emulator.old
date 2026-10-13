package tests

import (
	"nes-emu/src/emulator/ppu/internal/domain"
	"nes-emu/src/emulator/ppu/internal/strategies"
	shared "nes-emu/src/emulator/shared/application"
	"nes-emu/test/fixtures"
	"testing"
)

func TestCopyScrollXAndYToScrollRegister(t *testing.T) {
	sut, _ := createSut()

	if sut == nil {
		t.Error("Failed to create CopyScrollXAndYToScrollRegister strategy")
	}
}

func TestCopyScrollXAndYToScrollRegister_PanicWhenNotInitialized(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when strategy is not initialized, but did not panic")
		}
	}()

	sut := &strategies.CopyScrollToVRegister{}
	sut.Handle(nil)
}

func TestCopyScrollXAndYToScrollRegister_Handle(t *testing.T) {
	sut, bus := createSut()
	state := domain.NewPPU()

	bus.WriteToMemory(0x2005, 0x3F)
	sut.Handle(state)

	bus.WriteToMemory(0x2005, 0x7F)
	sut.Handle(state)

	if state.GetV() != 0x3F7F {
		t.Errorf("Expected V to be 0x3F7F, but got 0x%04X", state.GetV())
	}
}

func createSut() (*strategies.CopyScrollToVRegister, shared.MNIBus) {
	bus := fixtures.NewMockBus()
	return strategies.NewCopyScrollToVRegister(bus), bus
}
