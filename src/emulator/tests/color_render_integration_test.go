package tests

import (
	"fmt"
	"nes-emu/src/emulator/cpu"
	"nes-emu/src/emulator/ppu"
	"nes-emu/src/emulator/ppu/services"
	"nes-emu/src/emulator/shared/application"
	emu_application "nes-emu/src/emulator/shared/application"
	shared_persistence "nes-emu/src/emulator/shared/persistance"
	shared_services "nes-emu/src/emulator/shared/services"
	"testing"
)

func Test_RenderBackgroundColor(t *testing.T) {
	bus, cpu, ppu, screen := instantiateComponents()

	var program map[uint16]byte

	program = map[uint16]byte{
		//LDA $2004
		0x8000: 0xAD,
		0x8001: 0x04,
		0x8002: 0x20,

		//LDA #3F
		0x8003: 0xA9,
		0x8004: 0x3F,

		//STA $2006
		0x8005: 0x8D,
		0x8006: 0x06,
		0x8007: 0x20,

		//LDA #00
		0x8008: 0xA9,
		0x8009: 0x0,

		//STA $2006
		0x800A: 0x8D,
		0x800B: 0x06,
		0x800C: 0x20,

		//LDA #00
		0x800D: 0xA9,
		0x800E: 0x01,

		//STA $2007
		0x800F: 0x8D,
		0x8010: 0x07,
		0x8011: 0x20,

		//LDA #80
		0x8012: 0xA9,
		0x8013: 128,

		//STA $2000
		0x8014: 0x8D,
		0x8015: 0x00,
		0x8016: 0x20,

		//NOP
		0x8017: 0xEA,

		//JMP $8017
		0x8018: 0x4C,
		0x8019: 0x17,
		0x801A: 0x80,

		//NOP
		0x8100: 0xEA,
		0x8101: 0xEA,

		//NMI vector
		0xFFFA: 0x00,
		0xFFFB: 0x81,

		//Init vector
		0xFFFC: 0x00,
		0xFFFD: 0x80,
	}

	for addr, byte := range program {
		bus.WriteToMemory(addr, byte)
	}

	cpu.InitProgramCounter()

	if cpu.GetProgramCounter() != 0x8000 {
		t.Error("Program counter is not initialized")
	}

	cpu.StepIntoNextInstruction()

	if cpu.GetProgramCounter() != 0x8003 {
		t.Error("Program counter is not working")
	}

	cpu.StepIntoNextInstruction() //LDA
	cpu.StepIntoNextInstruction() //STA
	cpu.StepIntoNextInstruction() //LDA
	cpu.StepIntoNextInstruction() //STA

	if cpu.GetProgramCounter() != 0x800D {
		t.Errorf("Program counter is not working %X", cpu.GetProgramCounter())
	}

	if ppu.GetV() != 0x3F00 {
		t.Errorf("V is not working %X", ppu.GetV())
	}

	cpu.StepIntoNextInstruction() //LDA #80
	cpu.StepIntoNextInstruction() //STA $2000

	cycles := 0
	for cpu.GetProgramCounter() <= 0x8100 {
		cpu.StepIntoNextInstruction()
		cycles++

		fmt.Printf("%X", cpu.GetProgramCounter())

		if cycles > 100000 {
			t.Fatal("infinite loop detected")
		}
	}

	if cpu.GetProgramCounter() >= 0x8100 {
		fmt.Print("")
	}

	cpu.StepIntoNextInstruction()

	buffer := screen.GetImageBuffer()

	fmt.Printf("%v", buffer)
}

func instantiateComponents() (emu_application.Bus, TestCPU, TestPPU, *services.MemoryScreen) {
	bus := shared_services.NewBus()
	cpu := cpu.NewCPUModule(bus).GetController()
	ppuModule, screen := ppu.NewPPUModuleWithExternalScreenForTests(bus)
	ppu := ppuModule.GetController()

	bus.AtatchWorkMemory(shared_persistence.NewMemory())
	bus.AtatchVideoMemory(shared_persistence.NewMemory())
	bus.AttachNMI(cpu)
	bus.AttachPictureProcessingUnit(ppu)

	return bus, cpu, ppu, screen
}

type TestCPU interface {
	emu_application.CPU
	InitProgramCounter()
	StepIntoNextInstruction()
	GetProgramCounter() uint16
}

type TestPPU interface {
	application.PPU
	GetV() uint16
}
