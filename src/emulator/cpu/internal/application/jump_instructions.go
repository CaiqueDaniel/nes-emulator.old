package application

func (c *CPU) JumpProgramCounterToValue(value uint16) {
	c.programCounter = value
}

func (c *CPU) JumpProgramCounterByIndirectValue(address uint16) {
	highPtrAddress := address & 0xFF00
	lowPtrAddress := byte(address&0x00FF) + 1
	lowAddress := uint16(c.readFromMemory(address))
	highAddress := uint16(c.readFromMemory(highPtrAddress | uint16(lowPtrAddress)))
	finalAddress := highAddress<<8 | lowAddress

	c.JumpProgramCounterToValue(finalAddress)
}

func (c *CPU) JumpProgramCounterToSubRoutine(value uint16) {
	result := c.programCounter - 1
	lowAddress := uint8(result)
	highAddress := uint8((result & 0xFF00) >> 8)

	c.PushValueToStack(highAddress)
	c.PushValueToStack(lowAddress)

	c.programCounter = value
}

func (c *CPU) ReturnFromSubRoutine() {
	lowAddress := c.PullValueFromStack()
	highAddress := c.PullValueFromStack()

	c.programCounter = uint16(lowAddress) + uint16(highAddress)<<8
	c.programCounter++
}

func (c *CPU) ReturnFromInterrupt() {
	flags := c.PullValueFromStack()
	lowAddress := c.PullValueFromStack()
	highAddress := c.PullValueFromStack()

	c.carry = flags&0b00000001 != 0
	c.zero = flags&0b00000010 != 0
	c.interrupt = flags&0b00000100 != 0
	c.decimal = flags&0b00001000 != 0
	c.overflow = flags&0b01000000 != 0
	c.negative = flags&0b10000000 != 0
	c.programCounter = uint16(lowAddress) + uint16(highAddress)<<8
}

func (c *CPU) Break() {
	pcHighAddress := uint8(c.programCounter >> 8)
	pcLowAddress := uint8(c.programCounter)

	c.interrupt = true
	c.bFlag = true

	c.PushValueToStack(pcHighAddress)
	c.PushValueToStack(pcLowAddress)
	c.PushFlagsIntoStack()

	newPcLowAddress := uint16(c.readFromMemory(0xFFFE))
	newPcHighAddress := uint16(c.readFromMemory(0xFFFF))

	c.programCounter = newPcLowAddress | (newPcHighAddress << 8)
}
