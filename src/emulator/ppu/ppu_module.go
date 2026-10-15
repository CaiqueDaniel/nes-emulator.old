package ppu

import (
	"nes-emu/src/emulator/ppu/delivery"
	"nes-emu/src/emulator/ppu/internal/application"
	"nes-emu/src/emulator/ppu/internal/persistence"
	"nes-emu/src/emulator/ppu/internal/services"
	"nes-emu/src/emulator/ppu/internal/strategies"
	mock_services "nes-emu/src/emulator/ppu/services"
	shared_application "nes-emu/src/emulator/shared/application"

	"golang.org/x/exp/shiny/screen"
)

type PPUModule struct {
	initilized bool
	controller *delivery.PPUController
}

func NewPPUModule(bus shared_application.MNIBus, window *screen.Window, buffer *screen.Buffer) *PPUModule {
	pipeline := application.NewPipeline(bus)
	screen := services.NewShinyScreen(window, buffer)
	repository := persistence.NewPPUMemoryRepository()
	ioEventsContext := strategies.NewPPUIOEventContext(bus)
	renderGraphics := application.NewRenderGraphics(bus, pipeline, screen, repository)
	updateStateFromIOEvents := application.NewUpdateStateFromIOEvents(repository, ioEventsContext)

	controller := delivery.NewPPUController(renderGraphics, updateStateFromIOEvents)

	return &PPUModule{
		initilized: true,
		controller: controller,
	}
}

func NewPPUModuleWithExternalScreenForTests(bus shared_application.MNIBus) (*PPUModule, *mock_services.MemoryScreen) {
	pipeline := application.NewPipeline(bus)
	memoryScreen := mock_services.NewMemoryScreen()
	repository := persistence.NewPPUMemoryRepository()
	ioEventsContext := strategies.NewPPUIOEventContext(bus)
	renderGraphics := application.NewRenderGraphics(bus, pipeline, memoryScreen, repository)
	updateStateFromIOEvents := application.NewUpdateStateFromIOEvents(repository, ioEventsContext)

	controller := delivery.NewPPUController(renderGraphics, updateStateFromIOEvents)

	return &PPUModule{
		initilized: true,
		controller: controller,
	}, memoryScreen
}

func (p *PPUModule) GetController() *delivery.PPUController {
	return p.controller
}

func (p *PPUModule) checkIfInitilized() {
	if !p.initilized {
		panic("PPUModule was not initilized")
	}
}
