package application

import (
	"nes-emu/src/emulator/ppu/internal/domain"
	shared "nes-emu/src/emulator/shared/application"
)

type RenderGraphics struct {
	initialized bool
	bus         shared.MNIBus
	pipeline    PixelPipeline
	screen      Screen
	repository  domain.PPURepository
	buffer      [domain.MAX_FRAME_SCANLINE + 1][domain.MAX_PIXEL_PER_SCANLINE + 1]uint32
}

type RenderGraphicsInput struct {
	Address uint16
	IsWrite bool
}

func NewRenderGraphics(bus shared.MNIBus, pipeline PixelPipeline, screen Screen, repository domain.PPURepository) *RenderGraphics {
	p := &RenderGraphics{
		initialized: true,
		pipeline:    pipeline,
		bus:         bus,
		screen:      screen,
		buffer:      [domain.MAX_FRAME_SCANLINE + 1][domain.MAX_PIXEL_PER_SCANLINE + 1]uint32{},
		repository:  repository,
	}

	return p
}

func (p *RenderGraphics) Execute(input *RenderGraphicsInput) {
	p.checkIfInitilized()

	state := p.repository.GetState()

	if state.ShouldRenderScanlinePixel() {
		p.renderPixel(&state)
		state.AdvanceToNextScanlinePixel()
	}

	state.AdvanceToNextScanline()
	p.fetchGraphics(&state)
	p.updateStatusRegister(&state)

	if state.IsVBlankStarted() {
		p.screen.ShowImage(&p.buffer)
	}

	if p.checkIfNMIShouldBeCalled(&state) {
		p.bus.CallNMIHandler()
	}

	state.IncreaseDots()
	p.repository.Save(state)
}

func (p *RenderGraphics) GetCurrentScanline() uint16 {
	state := p.repository.GetState()
	return state.GetCurrentScanline()
}

func (p *RenderGraphics) GetCurrentScanlinePixel() uint8 {
	state := p.repository.GetState()
	return state.GetCurrentScanlinePixel()
}

func (p *RenderGraphics) GetV() uint16 {
	state := p.repository.GetState()
	return state.GetV()
}

func (p *RenderGraphics) renderPixel(state *domain.PPU) {
	pixelColor := p.pipeline.RenderPixel(state.GetX())

	p.buffer[state.GetCurrentScanline()][state.GetCurrentScanlinePixel()] = pixelColor
}

func (p *RenderGraphics) fetchGraphics(state *domain.PPU) {
	p.pipeline.StepUpPipeline(uint(state.GetCurrentDots()), state.GetV(), state.GetFineY())
}

func (p *RenderGraphics) updateStatusRegister(state *domain.PPU) {
	if state.IsOnPreRender() {
		value := p.bus.ReadFromMemory(domain.PPU_STATUS) & 0b00011111
		p.bus.WriteToMemory(domain.PPU_STATUS, value)
	}

	if state.IsVBlankStarted() {
		value := p.bus.ReadFromMemory(domain.PPU_STATUS) ^ 0b10000000
		p.bus.WriteToMemory(domain.PPU_STATUS, value)
	}
}

func (p *RenderGraphics) checkIfNMIShouldBeCalled(state *domain.PPU) bool {
	return p.isNMIFlagEnabled() && state.IsVBlankStarted()
}

func (p *RenderGraphics) isNMIFlagEnabled() bool {
	const nmiFlagMask = 0b10000000
	return p.bus.ReadFromMemory(domain.PPU_CONTROL)&nmiFlagMask != 0
}

func (p *RenderGraphics) checkIfInitilized() {
	if !p.initialized {
		panic("RenderGraphics was not initialized")
	}
}

func (p *RenderGraphics) readVMemory(address uint16) uint8 {
	return p.bus.ReadFromVideoMemory(address)
}

func (p *RenderGraphics) writeToVMemory(address uint16, value byte) {
	p.bus.WriteToVideoMemory(address, value)
}
