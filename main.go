package main

import (
	"image"
	"log"
	"nes-emu/src/emulator"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
)

func main() {
	driver.Main(func(s screen.Screen) {
		width, height := 260, 240
		size := image.Point{X: width, Y: height}

		// 1. Cria a janela
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title:  "Exemplo de Buffer de Pixels - Shiny",
			Width:  width,
			Height: height,
		})
		if err != nil {
			log.Fatalf("Erro ao criar janela: %v", err)
		}
		defer w.Release()

		// 2. Aloca um buffer de pixels na memória RAM
		buf, err := s.NewBuffer(size)
		if err != nil {
			log.Fatalf("Erro ao alocar buffer: %v", err)
		}
		defer buf.Release()

		emulator := emulator.NewEmulatorModule(&w, &buf)

		// 2. Goroutine para injetar eventos de renderização na fila de eventos da janela
		go func() {
			defer func() {
				if r := recover(); r != nil {
					panic(r)
				}
			}()

			err := emulator.GetStartGameController().StartGame("./test/resources/no-title")

			if err != nil {
				panic(err)
			}
		}()

		// 4. Laço principal de eventos (Event Loop manual)
		for {
			event := w.NextEvent()

			switch e := event.(type) {

			case lifecycle.Event:
				// Fecha a aplicação quando a janela é encerrada
				if e.To == lifecycle.StageDead {
					return
				}
			}
		}
	})
}
