package tests

import (
	ppu "nes-emu/src/emulator/ppu/internal/application"
	"nes-emu/src/emulator/ppu/internal/persistence"
	memory "nes-emu/src/emulator/shared/persistance"
	bus "nes-emu/src/emulator/shared/services"
	"nes-emu/test/fixtures"
	"testing"
)

func TestNewRenderGraphics(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	repository := persistence.NewPPUMemoryRepository()
	ppu := ppu.NewRenderGraphics(bus, mockPipeline, screenFixture, repository)

	if ppu == nil {
		t.Fatal("expected NewRenderGraphics to return a non-nil instance")
	}

	if ppu.GetCurrentScanline() != 0 {
		t.Errorf("expected initial scanline to be 0, got %d", ppu.GetCurrentScanline())
	}

	if ppu.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected initial pixel to be 0, got %d", ppu.GetCurrentScanlinePixel())
	}
}

func TestPPURender_ShouldPanic_WhenNotInitialized(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic when RenderGraphics is not initialized")
		}
		if r != "RenderGraphics was not initialized" {
			t.Errorf("expected panic message 'RenderGraphics was not initialized', got '%v'", r)
		}
	}()

	sut := &ppu.RenderGraphics{}
	sut.Execute(nil)
}

func TestPPURender_ShouldDrawAPixel(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	repository := persistence.NewPPUMemoryRepository()
	ppu := ppu.NewRenderGraphics(bus, mockPipeline, screenFixture, repository)

	ppu.Execute(nil)

	if ppu.GetCurrentScanline() != 0 {
		t.Fatal("scanline expected to be 0")
	}

	if ppu.GetCurrentScanlinePixel() != 1 {
		t.Fatal("expected to draw a pixel")
	}
}

func TestPPURender_ShouldWrapScanlineToStart(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	repository := persistence.NewPPUMemoryRepository()
	ppu := ppu.NewRenderGraphics(bus, mockPipeline, screenFixture, repository)

	for i := 0; i < 336*262; i++ {
		ppu.Execute(nil)
	}

	if ppu.GetCurrentScanline() != 0 {
		t.Errorf("scanline expected to be 0, got %d", ppu.GetCurrentScanline())
	}

	if ppu.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected to be at pixel 0, got %d", ppu.GetCurrentScanlinePixel())
	}
}

func TestPPURender_ShouldNotTriggerAnNMIOnVBlank_WhenNMIFlagDisabled(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockCpu := &fixtures.MockCPU{}
	mockPipeline := &PixelPipelineFixture{}

	bus.AttachNMI(mockCpu)

	screenFixture := NewScreenFixture()
	repository := persistence.NewPPUMemoryRepository()
	ppu := ppu.NewRenderGraphics(bus, mockPipeline, screenFixture, repository)

	for i := 0; i < 336*240; i++ {
		ppu.Execute(nil)
	}

	if ppu.GetCurrentScanline() != 240 {
		t.Errorf("scanline expected to be 240, got %d", ppu.GetCurrentScanline())
	}

	if ppu.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected to be at pixel 0, got %d", ppu.GetCurrentScanlinePixel())
	}

	if mockCpu.NmiCalled > 0 {
		t.Errorf("did not expect NMI to be called")
	}
}

func TestPPURender_ShouldTriggerAnNMIOnVBlank_WhenNMIFlagEnabled(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockCpu := &fixtures.MockCPU{}
	mockPipeline := &PixelPipelineFixture{}

	bus.AttachNMI(mockCpu)

	screenFixture := NewScreenFixture()
	repository := persistence.NewPPUMemoryRepository()
	ppu := ppu.NewRenderGraphics(bus, mockPipeline, screenFixture, repository)

	mem.Write(0x2000, 0b10000000)

	for i := 0; i < 336*240; i++ {
		ppu.Execute(nil)
	}

	if ppu.GetCurrentScanline() != 240 {
		t.Errorf("scanline expected to be 240, got %d", ppu.GetCurrentScanline())
	}

	if ppu.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected to be at pixel 0, got %d", ppu.GetCurrentScanlinePixel())
	}

	if mockCpu.NmiCalled == 0 {
		t.Errorf("did expect NMI to be called")
	}
}

func TestPPURender_ShouldResetFlagsOnStatusRegister_OnPreRender(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockCpu := &fixtures.MockCPU{}
	mockPipeline := &PixelPipelineFixture{}

	bus.AttachNMI(mockCpu)

	screenFixture := NewScreenFixture()
	repository := persistence.NewPPUMemoryRepository()
	ppu := ppu.NewRenderGraphics(bus, mockPipeline, screenFixture, repository)

	for i := 0; i < 336*261; i++ {
		ppu.Execute(nil)
	}

	if ppu.GetCurrentScanline() != 261 {
		t.Errorf("scanline expected to be 261, got %d", ppu.GetCurrentScanline())
	}

	if ppu.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected to be at pixel 0, got %d", ppu.GetCurrentScanlinePixel())
	}

	if mem.Read(0x2002) != 0 {
		t.Errorf("did expect flags to be cleared")
	}
}

func TestPPURender_ShouldResetFlagsOnStatusRegister_WithoutChangingOtherBits_OnPreRender(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockCpu := &fixtures.MockCPU{}

	bus.AttachNMI(mockCpu)

	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	repository := persistence.NewPPUMemoryRepository()
	ppu := ppu.NewRenderGraphics(bus, mockPipeline, screenFixture, repository)

	mem.Write(0x2002, 0b1111_1111)

	for i := 0; i < 336*261; i++ {
		ppu.Execute(nil)
	}

	if ppu.GetCurrentScanline() != 261 {
		t.Errorf("scanline expected to be 0, got %d", ppu.GetCurrentScanline())
	}

	if ppu.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected to be at pixel 0, got %d", ppu.GetCurrentScanlinePixel())
	}

	if mem.Read(0x2002) != 0b0001_1111 {
		t.Errorf("did expect value to be 0b0001_1111. got %b", mem.Read(0x2002))
	}
}

func TestPPURender_ShouldSetVBlankFlagOnStatusRegister_OnVBlank(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockCpu := &fixtures.MockCPU{}
	mockPipeline := &PixelPipelineFixture{}

	bus.AttachNMI(mockCpu)

	screenFixture := NewScreenFixture()
	repository := persistence.NewPPUMemoryRepository()
	ppu := ppu.NewRenderGraphics(bus, mockPipeline, screenFixture, repository)

	for i := 0; i < 336*240; i++ {
		ppu.Execute(nil)
	}

	if ppu.GetCurrentScanline() != 240 {
		t.Errorf("scanline expected to be 0, got %d", ppu.GetCurrentScanline())
	}

	if ppu.GetCurrentScanlinePixel() != 0 {
		t.Errorf("expected to be at pixel 0, got %d", ppu.GetCurrentScanlinePixel())
	}

	if mem.Read(0x2002) != 0b1000_0000 {
		t.Errorf("did expect v-blank to be setted")
	}
}

func TestPPURender_ShouldShiftRegisters_OnVisibleScanlines(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	repository := persistence.NewPPUMemoryRepository()
	sut := ppu.NewRenderGraphics(bus, mockPipeline, screenFixture, repository)

	for i := 0; i <= 16; i++ {
		sut.Execute(nil)
	}

	lowPatternShiftRegister, highPatternShiftRegister, lowAttributeShiftRegister, highAttributeShiftRegister := mockPipeline.GetShiftRegisters()

	if lowPatternShiftRegister != 0b10_00000010 {
		t.Errorf("expected low pattern shift register to be 0b00000010_00000010, got %b", lowPatternShiftRegister)
	}

	if highPatternShiftRegister != 0b11_00000011 {
		t.Errorf("expected high pattern shift register to be 0b00000011_00000011, got %b", highPatternShiftRegister)
	}

	if lowAttributeShiftRegister != 0b1_00000001 {
		t.Errorf("expected low attribute shift register to be 0b00000001_00000001, got %b", lowAttributeShiftRegister)
	}

	if highAttributeShiftRegister != 0b1_00000001 {
		t.Errorf("expected high attribute shift register to be 0b00000001_00000001, got %b", highAttributeShiftRegister)
	}
}

func TestPPURender_ShouldShiftRegisters_OnHBlank(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	repository := persistence.NewPPUMemoryRepository()
	sut := ppu.NewRenderGraphics(bus, mockPipeline, screenFixture, repository)

	for i := 0; i <= 256; i++ {
		sut.Execute(nil)
	}

	lowPatternShiftRegister, highPatternShiftRegister, lowAttributeShiftRegister, highAttributeShiftRegister := mockPipeline.GetShiftRegisters()

	if lowPatternShiftRegister != 0b1000_00010 {
		t.Errorf("expected low pattern shift register to be 0b1000_00010, got %b", lowPatternShiftRegister)
	}

	if highPatternShiftRegister != 0b1000000110000011 {
		t.Errorf("expected high pattern shift register to be 0b1000000110000011, got %b", highPatternShiftRegister)
	}

	if lowAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected low attribute shift register to be 0b1000000010000001, got %b", lowAttributeShiftRegister)
	}

	if highAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected high attribute shift register to be 0b1000000010000001, got %b", highAttributeShiftRegister)
	}

	sut.Execute(nil)
	lowPatternShiftRegister, highPatternShiftRegister, lowAttributeShiftRegister, highAttributeShiftRegister = mockPipeline.GetShiftRegisters()

	if lowPatternShiftRegister != 0b1000_00010 {
		t.Errorf("expected low pattern shift register to be 0b1000_00010, got %b", lowPatternShiftRegister)
	}

	if highPatternShiftRegister != 0b1000000110000011 {
		t.Errorf("expected high pattern shift register to be 0b1000000110000011, got %b", highPatternShiftRegister)
	}

	if lowAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected low attribute shift register to be 0b1000000010000001, got %b", lowAttributeShiftRegister)
	}

	if highAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected high attribute shift register to be 0b1000000010000001, got %b", highAttributeShiftRegister)
	}
}

func TestPPURender_ShouldShiftRegisters_OnVBlank(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screen := NewScreenFixture()
	repository := persistence.NewPPUMemoryRepository()
	sut := ppu.NewRenderGraphics(bus, mockPipeline, screen, repository)

	for i := 0; i <= 336*240; i++ {
		sut.Execute(nil)
	}

	lowPatternShiftRegister, highPatternShiftRegister, lowAttributeShiftRegister, highAttributeShiftRegister := mockPipeline.GetShiftRegisters()

	if lowPatternShiftRegister != 0b100000010 {
		t.Errorf("expected low pattern shift register to be 0b100000010, got %b", lowPatternShiftRegister)
	}

	if highPatternShiftRegister != 0b1000000110000011 {
		t.Errorf("expected high pattern shift register to be 0b1000000110000011, got %b", highPatternShiftRegister)
	}

	if lowAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected low attribute shift register to be 0b1000000010000001, got %b", lowAttributeShiftRegister)
	}

	if highAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected high attribute shift register to be 0b1000000010000001, got %b", highAttributeShiftRegister)
	}

	sut.Execute(nil)
	lowPatternShiftRegister, highPatternShiftRegister, lowAttributeShiftRegister, highAttributeShiftRegister = mockPipeline.GetShiftRegisters()

	if lowPatternShiftRegister != 0b100000010 {
		t.Errorf("expected low pattern shift register to be 0b100000010, got %b", lowPatternShiftRegister)
	}

	if highPatternShiftRegister != 0b1000000110000011 {
		t.Errorf("expected high pattern shift register to be 0b1000000110000011, got %b", highPatternShiftRegister)
	}

	if lowAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected low attribute shift register to be 0b1000000010000001, got %b", lowAttributeShiftRegister)
	}

	if highAttributeShiftRegister != 0b1000000010000001 {
		t.Errorf("expected high attribute shift register to be 0b1000000010000001, got %b", highAttributeShiftRegister)
	}
}

func TestPPURender_ShouldShowImageOnScreen_OnVBlankStart(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	repository := persistence.NewPPUMemoryRepository()
	ppu := ppu.NewRenderGraphics(bus, mockPipeline, screenFixture, repository)

	// Render until start of VBlank (scanline 240)
	for i := 0; i < 336*240; i++ {
		ppu.Execute(nil)
	}

	if screenFixture.ShowImageCalls == 0 {
		t.Errorf("expected ShowImage to be called on VBlank start, got %d", screenFixture.ShowImageCalls)
	}

	if screenFixture.LastBuffer == nil {
		t.Error("expected LastBuffer to not be nil")
	}

	// Verify that visible scanlines in the buffer have rendered pixels (256 pixels each)
	buffer := *screenFixture.LastBuffer
	if len(buffer[0]) != 256 {
		t.Errorf("expected 256 pixels rendered for scanline 0, got %d", len(buffer[0]))
	}
}

func TestPPURender_ShouldNotRenderPixel_DuringHBlank(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	repository := persistence.NewPPUMemoryRepository()
	ppu := ppu.NewRenderGraphics(bus, mockPipeline, screenFixture, repository)

	// Execute through the visible portion of scanline 0 (dots 0 to 255: 256 cycles)
	for i := 0; i < 256; i++ {
		ppu.Execute(nil)
	}

	pixelAfterVisible := ppu.GetCurrentScanlinePixel()

	// Execute through HBlank portion of scanline 0 (dots 256 to 335: 80 cycles)
	for i := 0; i < 80; i++ {
		ppu.Execute(nil)
		if ppu.GetCurrentScanlinePixel() != pixelAfterVisible {
			t.Fatalf("expected pixel to remain unchanged during HBlank at dot %d, got %d", 256+i, ppu.GetCurrentScanlinePixel())
		}
	}
}

func TestPPURender_ShouldNotRenderPixel_DuringVBlank(t *testing.T) {
	mem := memory.NewMemory()
	vMemory := memory.NewMemory()
	bus := bus.NewBusWithWorkMemory(mem)
	bus.AtatchVideoMemory(vMemory)
	mockPipeline := &PixelPipelineFixture{}
	screenFixture := NewScreenFixture()
	repository := persistence.NewPPUMemoryRepository()
	ppu := ppu.NewRenderGraphics(bus, mockPipeline, screenFixture, repository)

	// Execute until VBlank begins (scanline 240, dot 0)
	for i := 0; i < 336*240; i++ {
		ppu.Execute(nil)
	}

	if ppu.GetCurrentScanline() != 240 {
		t.Fatalf("expected scanline 240, got %d", ppu.GetCurrentScanline())
	}

	pixelAtVBlankStart := ppu.GetCurrentScanlinePixel()

	// Execute 336 cycles across scanline 240 during VBlank
	for i := 0; i < 336; i++ {
		ppu.Execute(nil)
		if ppu.GetCurrentScanlinePixel() != pixelAtVBlankStart {
			t.Fatalf("expected pixel to not advance during VBlank scanline at step %d, got %d", i, ppu.GetCurrentScanlinePixel())
		}
	}
}
