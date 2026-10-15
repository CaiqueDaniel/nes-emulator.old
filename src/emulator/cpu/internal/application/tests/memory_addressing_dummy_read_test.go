package tests

import (
	"reflect"
	"testing"

	internal "nes-emu/src/emulator/cpu/internal/application"
	shared "nes-emu/src/emulator/shared/application"
	memory "nes-emu/src/emulator/shared/persistance"
	bus "nes-emu/src/emulator/shared/services"
)

type addressRecordingBus struct {
	shared.Bus
	readAddresses []uint16
}

func (b *addressRecordingBus) ReadFromMemory(address uint16) uint8 {
	b.readAddresses = append(b.readAddresses, address)
	return b.Bus.ReadFromMemory(address)
}

func newCPUWithAddressRecordingBus() (*internal.CPU, shared.Memory, *addressRecordingBus) {
	mem := memory.NewMemory()
	recordingBus := &addressRecordingBus{
		Bus: bus.NewBusWithWorkMemory(mem),
	}
	return internal.NewCpuWithInternal(recordingBus), mem, recordingBus
}

func TestGetAddressByIndexedAbsoluteModeDummyReadOnPageCross(t *testing.T) {
	cpu, _, recordingBus := newCPUWithAddressRecordingBus()

	address := cpu.GetAddressByIndexedAbsoluteMode(0x12FF, 0x01)

	if address != 0x1300 {
		t.Fatalf("Expected effective address 0x1300, got 0x%04X", address)
	}
	if !reflect.DeepEqual(recordingBus.readAddresses, []uint16{0x1200}) {
		t.Errorf("Expected dummy read at 0x1200, got %04X", recordingBus.readAddresses)
	}
}

func TestGetAddressByIndexedAbsoluteModeWithDummyReadReadsProvisionalAddressOnce(t *testing.T) {
	cpu, _, recordingBus := newCPUWithAddressRecordingBus()

	address := cpu.GetAddressByIndexedAbsoluteModeWithDummyRead(0x12FF, 0x01)

	if address != 0x1300 {
		t.Fatalf("Expected effective address 0x1300, got 0x%04X", address)
	}
	if !reflect.DeepEqual(recordingBus.readAddresses, []uint16{0x1200}) {
		t.Errorf("Expected one dummy read at 0x1200, got %04X", recordingBus.readAddresses)
	}
}

func TestGetAddressByIndirectIndexedYModeDummyReadOnPageCross(t *testing.T) {
	cpu, mem, recordingBus := newCPUWithAddressRecordingBus()
	mem.Write(0x0010, 0xFF)
	mem.Write(0x0011, 0x12)
	cpu.LoadValueIntoRegister(0x01, internal.REGISTER_Y)

	address := cpu.GetAddressByIndirectIndexedYMode(0x10)

	if address != 0x1300 {
		t.Fatalf("Expected effective address 0x1300, got 0x%04X", address)
	}
	expectedReads := []uint16{0x0010, 0x0011, 0x1200}
	if !reflect.DeepEqual(recordingBus.readAddresses, expectedReads) {
		t.Errorf("Expected reads %04X, got %04X", expectedReads, recordingBus.readAddresses)
	}
}

func TestGetAddressByIndirectIndexedYModeWithDummyReadUsesProvisionalAddress(t *testing.T) {
	cpu, mem, recordingBus := newCPUWithAddressRecordingBus()
	mem.Write(0x0010, 0xFF)
	mem.Write(0x0011, 0x12)
	cpu.LoadValueIntoRegister(0x01, internal.REGISTER_Y)

	address := cpu.GetAddressByIndirectIndexedYModeWithDummyRead(0x10)

	if address != 0x1300 {
		t.Fatalf("Expected effective address 0x1300, got 0x%04X", address)
	}
	expectedReads := []uint16{0x0010, 0x0011, 0x1200}
	if !reflect.DeepEqual(recordingBus.readAddresses, expectedReads) {
		t.Errorf("Expected reads %04X, got %04X", expectedReads, recordingBus.readAddresses)
	}
}
