package application

import (
	"fmt"
	"nes-emu/src/emulator/shared/application"
	"slices"
	"time"
)

const REGISTER_X = "X"
const REGISTER_Y = "Y"
const ACCUMULATOR = "ACC"
const STACK_START = 0x01FF
const STACK_END = 0x0100
const START_POINTER = 0xFFFC
const clock_in_mhz = 1.789773

type CPU struct {
	programCounter                                                  uint16
	acc, x, y                                                       uint8
	carry, zero, overflow, negative, interrupt, decimal, bFlag, nmi bool
	stackPointer                                                    uint8
	memory                                                          application.Memory
	bus                                                             application.Bus
	stopPcAt                                                        int
	instructionSet                                                  instructionSet
	currentFrameCycles                                              uint16
	stopProgram                                                     bool
}

func NewCpu(bus application.Bus) *CPU {
	c := &CPU{
		bus:      bus,
		stopPcAt: -1,
	}
	c.instructionSet = c.initInstructions()
	c.Reset()
	return c
}

func NewCpuWithStopAt(bus application.Bus, stopPcAt int) *CPU {
	c := &CPU{
		bus:      bus,
		stopPcAt: stopPcAt,
	}
	c.instructionSet = c.initInstructions()
	c.Reset()
	return c
}

func NewCpuWithProgramCounter(programCounter uint16, bus application.Bus) *CPU {
	c := &CPU{
		programCounter: programCounter,
		stackPointer:   0xFF,
		bus:            bus,
		stopPcAt:       -1,
	}
	c.instructionSet = c.initInstructions()
	return c
}

func NewCpuWithInternal(bus application.Bus) *CPU {
	c := &CPU{
		bus:      bus,
		stopPcAt: -1,
	}
	c.instructionSet = c.initInstructions()
	c.Reset()
	return c
}

func (c *CPU) RunProgram() {
	c.InitProgramCounter()
	c.runGameLoop()
}

func (c *CPU) Reset() {
	c.programCounter = 0
	c.acc = 0
	c.x = 0
	c.y = 0
	c.carry = false
	c.zero = false
	c.overflow = false
	c.negative = false
	c.interrupt = true
	c.decimal = false
	c.bFlag = false
	c.stackPointer = 0xFF
	c.nmi = false
	c.currentFrameCycles = 0
	c.stopProgram = false
}

func (c *CPU) SetNMI() {
	c.nmi = true
}

func (c *CPU) HandleNMI() {
	const nmi_handler_byte_ptr = 0xFFFA

	highPcByte := byte((c.programCounter & 0xFF00) >> 8)
	lowPcByte := byte(c.programCounter)
	c.bFlag = false

	c.PushValueToStack(highPcByte)
	c.PushValueToStack(lowPcByte)
	c.PushFlagsIntoStack()
	c.SetInterruptFlag()

	lowByte := c.readFromMemory(nmi_handler_byte_ptr)
	highByte := c.readFromMemory(nmi_handler_byte_ptr + 1)

	c.programCounter = uint16(highByte)<<8 | uint16(lowByte)
	c.nmi = false
}

func (c *CPU) PushValueToStack(value uint8) {
	c.writeToMemory(STACK_END+uint16(c.stackPointer), value)
	c.stackPointer--
}

func (c *CPU) PullValueFromStack() uint8 {
	c.stackPointer++
	value := c.readFromMemory(STACK_END + uint16(c.stackPointer))

	c.writeToMemory(STACK_END+uint16(c.stackPointer), 0)

	return value
}

func (c *CPU) PushFlagsIntoStack() {
	carry := transformFlagIntoUint8(c.carry)
	zero := transformFlagIntoUint8(c.zero) << 1
	irq := transformFlagIntoUint8(c.interrupt) << 2
	decimal := transformFlagIntoUint8(c.decimal) << 3
	breakFlag := transformFlagIntoUint8(c.bFlag) << 4
	overflow := transformFlagIntoUint8(c.overflow) << 6
	negative := transformFlagIntoUint8(c.negative) << 7

	valueToPush := 0b00100000 | carry | zero | irq | decimal | breakFlag | overflow | negative
	c.PushValueToStack(valueToPush)
}

func (c *CPU) InitProgramCounter() {
	startAddressLow := c.readFromMemory(START_POINTER)
	startAddressHigh := c.readFromMemory(START_POINTER + 1)
	startAddress := (uint16(startAddressHigh) << 8) + uint16(startAddressLow)
	c.programCounter = startAddress
}

func (c *CPU) runGameLoop() {
	const time_per_frame = time.Second / 60
	startTime := time.Now()

	for {
		c.currentFrameCycles = 0

		c.renderFrame()

		timeElapsed := time.Since(startTime)

		if timeElapsed < time_per_frame {
			time.Sleep(time_per_frame - timeElapsed)
			startTime = time.Now()
		}

		if c.stopProgram {
			break
		}
	}
}

func (c *CPU) renderFrame() {
	const cycles_per_frame = 29781

	for c.currentFrameCycles < cycles_per_frame {
		c.RunInstruction()

		if c.stopPcAt != -1 && c.programCounter >= uint16(c.stopPcAt) {
			c.stopProgram = true
			break
		}
	}
}

func (c *CPU) RunInstruction() {
	if c.nmi {
		c.HandleNMI()
	}

	if c.programCounter == 0x8100 {
		fmt.Print()
	}

	opCode := c.readFromMemory(c.programCounter)
	c.programCounter++

	c.interpretInstruction(opCode)
}

func (c *CPU) interpretInstruction(opCode uint8) {
	instruction := c.instructionSet[opCode]

	if instruction == nil {
		return
	}

	bytes := make([]uint8, 0)

	for i := 0; i < instruction.ArgsBytes; i++ {
		bytes = append(bytes, c.readFromMemory(c.programCounter))
		c.programCounter++
	}

	if instruction.ArgsBytes == 0 {
		c.doDummyMemoryRead(c.programCounter)
	}

	slices.Reverse(bytes)
	instruction.Method(bytes)
}

func (c *CPU) writeToMemory(address uint16, value uint8) {
	c.bus.Tick()
	c.bus.WriteToMemory(address, value)
	c.bus.TriggerIOEvents(address, true)
	c.currentFrameCycles++
}

func (c *CPU) readFromMemory(address uint16) uint8 {
	c.bus.Tick()
	value := c.bus.ReadFromMemory(address)
	c.bus.TriggerIOEvents(address, false)
	c.currentFrameCycles++
	return value
}

func (c *CPU) doDummyMemoryRead(address uint16) {
	c.readFromMemory(address)
}

func transformFlagIntoUint8(flag bool) uint8 {
	if flag {
		return 1
	}
	return 0
}
