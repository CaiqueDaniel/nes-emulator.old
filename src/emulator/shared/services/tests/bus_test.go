package tests

import (
	"nes-emu/src/emulator/shared/application"
	memory "nes-emu/src/emulator/shared/persistance"
	bus "nes-emu/src/emulator/shared/services"
	"nes-emu/test/fixtures"
	"testing"
)

type mockPPU struct {
	renderCalled int
	events       []*application.PPUIOEvent
}

func (m *mockPPU) Render(request *application.PPUIOEvent) {
	m.renderCalled++
	if request != nil {
		eventCopy := *request
		m.events = append(m.events, &eventCopy)
	} else {
		m.events = append(m.events, nil)
	}
}

func (m *mockPPU) HandleIOEvents(address uint16, isWrite bool) {}

func TestNewBus(t *testing.T) {
	b := bus.NewBus()
	if b == nil {
		t.Fatalf("Expected NewBus to return a non-nil object")
	}

	if b.GetTickCount() != 0 {
		t.Errorf("Expected initial tickCount to be 0, got %d", b.GetTickCount())
	}
}

func TestNewBusWithWorkMemory(t *testing.T) {
	mem := memory.NewMemory()
	b := bus.NewBusWithWorkMemory(mem)
	if b == nil {
		t.Fatalf("Expected NewBusWithWorkMemory to return a non-nil object")
	}

	b.WriteToMemory(0x1234, 0x42)
	if b.ReadFromMemory(0x1234) != 0x42 {
		t.Errorf("Expected work memory to be initialized and accessible")
	}
}

func TestAttachWorkMemory(t *testing.T) {
	b := bus.NewBus()
	mem := memory.NewMemory()
	b.AtatchWorkMemory(mem)

	b.WriteToMemory(0x0100, 0xAB)
	result := b.ReadFromMemory(0x0100)
	if result != 0xAB {
		t.Errorf("Expected to read 0xAB from attached work memory, got 0x%X", result)
	}
}

func TestAttachVideoMemory(t *testing.T) {
	b := bus.NewBus()
	vmem := memory.NewMemory()
	b.AtatchVideoMemory(vmem)

	b.WriteToVideoMemory(0x2000, 0xCD)
	result := b.ReadFromVideoMemory(0x2000)
	if result != 0xCD {
		t.Errorf("Expected to read 0xCD from attached video memory, got 0x%X", result)
	}
}

func TestTick_WithoutPPU(t *testing.T) {
	b := bus.NewBus()

	b.Tick()
	if b.GetTickCount() != 1 {
		t.Errorf("Expected tickCount to be 1, got %d", b.GetTickCount())
	}
}

func TestAttachPictureProcessingUnitAndTick(t *testing.T) {
	b := bus.NewBus()
	ppu := &mockPPU{}
	b.AttachPictureProcessingUnit(ppu)

	b.Tick()

	if ppu.renderCalled != 3 {
		t.Errorf("Expected ppu.Render to be called 3 times after 1 tick, got %d", ppu.renderCalled)
	}

	if len(ppu.events) != 3 {
		t.Fatalf("Expected 3 render event records, got %d", len(ppu.events))
	}

	if ppu.events[1] != nil || ppu.events[2] != nil {
		t.Errorf("Expected subsequent render calls to receive nil event")
	}

	if b.GetTickCount() != 1 {
		t.Errorf("Expected tickCount to be 1, got %d", b.GetTickCount())
	}

	b.Tick()

	if ppu.renderCalled != 6 {
		t.Errorf("Expected ppu.Render to be called 6 times after 2 ticks, got %d", ppu.renderCalled)
	}

	if b.GetTickCount() != 2 {
		t.Errorf("Expected tickCount to be 2, got %d", b.GetTickCount())
	}
}

func TestResetTickCount(t *testing.T) {
	b := bus.NewBus()
	ppu := &mockPPU{}
	b.AttachPictureProcessingUnit(ppu)

	b.Tick()
	b.Tick()
	if b.GetTickCount() != 2 {
		t.Errorf("Expected tickCount to be 2, got %d", b.GetTickCount())
	}

	b.ResetTickCount()
	if b.GetTickCount() != 0 {
		t.Errorf("Expected tickCount to be reset to 0, got %d", b.GetTickCount())
	}
}

func TestAttachNMIAndCallNMIHandler(t *testing.T) {
	b := bus.NewBus()
	cpu := &fixtures.MockCPU{}

	b.AttachNMI(cpu)
	b.CallNMIHandler()

	if cpu.NmiCalled != 1 {
		t.Errorf("Expected cpu.SetNMI to be called 1 time, got %d", cpu.NmiCalled)
	}

	b.CallNMIHandler()

	if cpu.NmiCalled != 2 {
		t.Errorf("Expected cpu.SetNMI to be called 2 times, got %d", cpu.NmiCalled)
	}
}

func TestReadWorkMemory_PanicsWhenUnattached(t *testing.T) {
	b := bus.NewBus()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected ReadFromMemory to panic when work memory is nil")
		}
	}()
	b.ReadFromMemory(0x0000)
}

func TestWriteWorkMemory_PanicsWhenUnattached(t *testing.T) {
	b := bus.NewBus()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected WriteToMemory to panic when work memory is nil")
		}
	}()
	b.WriteToMemory(0x0000, 0xFF)
}

func TestReadVideoMemory_PanicsWhenUnattached(t *testing.T) {
	b := bus.NewBus()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected ReadFromVideoMemory to panic when video memory is nil")
		}
	}()
	b.ReadFromVideoMemory(0x0000)
}

func TestWriteVideoMemory_PanicsWhenUnattached(t *testing.T) {
	b := bus.NewBus()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected WriteToVideoMemory to panic when video memory is nil")
		}
	}()
	b.WriteToVideoMemory(0x0000, 0xFF)
}

func TestReadAndWriteToMemory_MemoryMirroing_RAM(t *testing.T) {
	b := bus.NewBus()
	b.AtatchWorkMemory(memory.NewMemory())

	b.WriteToMemory(0x0000, 0x12)
	b.WriteToMemory(0x7FF, 0x34)

	if b.ReadFromMemory(0x0000) != 0x12 || b.ReadFromMemory(0x800) != 0x12 || b.ReadFromMemory(0x1000) != 0x12 || b.ReadFromMemory(0x1800) != 0x12 {
		t.Errorf("Expected to read 0x12 from memory")
	}

	if b.ReadFromMemory(0x7FF) != 0x34 && b.ReadFromMemory(0x1FF) != 0x34 && b.ReadFromMemory(0x9FF) != 0x34 && b.ReadFromMemory(0x1FFF) != 0x34 {
		t.Errorf("Expected to read 0x34 from memory")
	}

	b.WriteToMemory(0x800, 0xAB)

	if b.ReadFromMemory(0x0000) != 0xAB {
		t.Errorf("Expected to read 0xAB from memory at address 0x0000, got 0x%X", b.ReadFromMemory(0x0000))
	}
}

func TestReadAndWriteToMemory_MemoryMirroing_PPULatches(t *testing.T) {
	b := bus.NewBus()
	b.AtatchWorkMemory(memory.NewMemory())

	b.WriteToMemory(0x2000, 0x12)
	b.WriteToMemory(0x2007, 0x34)

	if b.ReadFromMemory(0x2000) != 0x12 {
		t.Errorf("Expected to read 0x12 from memory")
	}

	if b.ReadFromMemory(0x2008) != 0x12 {
		t.Errorf("Expected to read 0x12 from memory")
	}

	if b.ReadFromMemory(0x2010) != 0x12 {
		t.Errorf("Expected to read 0x12 from memory")
	}

	if b.ReadFromMemory(0x2018) != 0x12 {
		t.Errorf("Expected to read 0x12 from memory")
	}

	if b.ReadFromMemory(0x2007) != 0x34 && b.ReadFromMemory(0x3FFF) != 0x34 {
		t.Errorf("Expected to read 0x34 from memory")
	}

	b.WriteToMemory(0x2008, 0xAB)

	if b.ReadFromMemory(0x2000) != 0xAB {
		t.Errorf("Expected to read 0xAB from memory at address 0x2000, got 0x%X", b.ReadFromMemory(0x2000))
	}
}

func TestReadAndWriteToMemory_NoMirroring_AbovePPULatches(t *testing.T) {
	b := bus.NewBus()
	b.AtatchWorkMemory(memory.NewMemory())

	// Addresses >= 0x4000 (APU, I/O, ROM) should not be mirrored
	b.WriteToMemory(0x4000, 0x11)
	b.WriteToMemory(0x4008, 0x22)
	b.WriteToMemory(0x8000, 0x33)
	b.WriteToMemory(0x8008, 0x44)

	if b.ReadFromMemory(0x4000) != 0x11 {
		t.Errorf("Expected 0x11 at 0x4000, got 0x%X", b.ReadFromMemory(0x4000))
	}
	if b.ReadFromMemory(0x4008) != 0x22 {
		t.Errorf("Expected 0x22 at 0x4008, got 0x%X", b.ReadFromMemory(0x4008))
	}
	if b.ReadFromMemory(0x8000) != 0x33 {
		t.Errorf("Expected 0x33 at 0x8000, got 0x%X", b.ReadFromMemory(0x8000))
	}
	if b.ReadFromMemory(0x8008) != 0x44 {
		t.Errorf("Expected 0x44 at 0x8008, got 0x%X", b.ReadFromMemory(0x8008))
	}
}

func TestReadAndWriteToVideoMemory_MemoryMirroring_Palettes(t *testing.T) {
	b := bus.NewBus()
	b.AtatchVideoMemory(memory.NewMemory())

	// Palettes are at 0x3F00 - 0x3F1F (32 bytes) and mirrored up to 0x3FFF
	b.WriteToVideoMemory(0x3F00, 0x55)
	b.WriteToVideoMemory(0x3F1F, 0x66)

	// Check mirroring of 0x3F00 at 0x3F20, 0x3F40, 0x3FE0
	if val := b.ReadFromVideoMemory(0x3F00); val != 0x55 {
		t.Errorf("Expected 0x55 at 0x3F00, got 0x%X", val)
	}
	if val := b.ReadFromVideoMemory(0x3F20); val != 0x55 {
		t.Errorf("Expected 0x55 at 0x3F20, got 0x%X", val)
	}
	if val := b.ReadFromVideoMemory(0x3F40); val != 0x55 {
		t.Errorf("Expected 0x55 at 0x3F40, got 0x%X", val)
	}
	if val := b.ReadFromVideoMemory(0x3FE0); val != 0x55 {
		t.Errorf("Expected 0x55 at 0x3FE0, got 0x%X", val)
	}

	// Check mirroring of 0x3F1F at 0x3F3F, 0x3FFF
	if val := b.ReadFromVideoMemory(0x3F1F); val != 0x66 {
		t.Errorf("Expected 0x66 at 0x3F1F, got 0x%X", val)
	}
	if val := b.ReadFromVideoMemory(0x3F3F); val != 0x66 {
		t.Errorf("Expected 0x66 at 0x3F3F, got 0x%X", val)
	}
	if val := b.ReadFromVideoMemory(0x3FFF); val != 0x66 {
		t.Errorf("Expected 0x66 at 0x3FFF, got 0x%X", val)
	}

	// Write to mirror should update base address
	b.WriteToVideoMemory(0x3F20, 0xAA)
	if val := b.ReadFromVideoMemory(0x3F00); val != 0xAA {
		t.Errorf("Expected 0xAA at 0x3F00 after write to mirror 0x3F20, got 0x%X", val)
	}

	b.WriteToVideoMemory(0x3FFF, 0xBB)
	if val := b.ReadFromVideoMemory(0x3F1F); val != 0xBB {
		t.Errorf("Expected 0xBB at 0x3F1F after write to mirror 0x3FFF, got 0x%X", val)
	}
}

func TestReadAndWriteToVideoMemory_NoMirroring_BelowPalettes(t *testing.T) {
	b := bus.NewBus()
	b.AtatchVideoMemory(memory.NewMemory())

	// Addresses below 0x3F00 (pattern tables, nametables) should not be mirrored by palette mirroring
	b.WriteToVideoMemory(0x0000, 0x10)
	b.WriteToVideoMemory(0x0020, 0x20)
	b.WriteToVideoMemory(0x2000, 0x30)
	b.WriteToVideoMemory(0x3EFF, 0x40)

	if val := b.ReadFromVideoMemory(0x0000); val != 0x10 {
		t.Errorf("Expected 0x10 at 0x0000, got 0x%X", val)
	}
	if val := b.ReadFromVideoMemory(0x0020); val != 0x20 {
		t.Errorf("Expected 0x20 at 0x0020, got 0x%X", val)
	}
	if val := b.ReadFromVideoMemory(0x2000); val != 0x30 {
		t.Errorf("Expected 0x30 at 0x2000, got 0x%X", val)
	}
	if val := b.ReadFromVideoMemory(0x3EFF); val != 0x40 {
		t.Errorf("Expected 0x40 at 0x3EFF, got 0x%X", val)
	}
}
