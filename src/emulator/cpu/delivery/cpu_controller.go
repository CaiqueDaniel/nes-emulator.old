package delivery

import (
	"fmt"
	"nes-emu/src/emulator/cpu/internal/application"
)

type CPUController struct {
	initialize bool
	cpu        *application.CPU
}

func NewCPUController(cpu *application.CPU) *CPUController {
	return &CPUController{
		initialize: true,
		cpu:        cpu,
	}
}

func (c *CPUController) RunProgram() {
	c.checkIfInitialized()
	c.cpu.RunProgram()
}

func (c *CPUController) Reset() {
	c.checkIfInitialized()
	c.cpu.Reset()
}

func (c *CPUController) SetNMI() {
	c.checkIfInitialized()
	c.cpu.SetNMI()
}

func (c *CPUController) InitProgramCounter() {
	c.checkIfInitialized()
	c.cpu.InitProgramCounter()
}

func (c *CPUController) StepIntoNextInstruction() {
	c.checkIfInitialized()
	c.cpu.RunInstruction()
}

func (c *CPUController) GetProgramCounter() uint16 {
	c.checkIfInitialized()
	a := c.cpu.GetProgramCounter()

	if a >= 0x8100 {
		fmt.Print()
	}

	return a
}

func (c *CPUController) checkIfInitialized() {
	if !c.initialize {
		panic("CPU not initialized")
	}
}
