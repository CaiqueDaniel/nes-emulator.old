package application

func (c *CPU) GetValueByAbsoluteMode(address uint16) uint8 {
	return c.readFromMemory(address)
}

func (c *CPU) GetValueByZeroPageMode(address uint8) uint8 {
	return c.readFromMemory(uint16(address))
}

func (c *CPU) GetValueByIndexedAbsoluteMode(address uint16, index uint8) uint8 {
	return c.readFromMemory(c.GetAddressByIndexedAbsoluteMode(address, index))
}

func (c *CPU) GetValueByZeroPageIndexedMode(address uint8, index uint8) uint8 {
	return c.readFromMemory(uint16(address + index))
}

func (c *CPU) GetValueByZeroPageIndexedModeWithDummyRead(address uint8, index uint8) uint8 {
	c.doDummyMemoryRead(uint16(address))
	return c.readFromMemory(uint16(address + index))
}

func (c *CPU) GetAddressByZeroPageIndexedModeWithDummyRead(address uint8, index uint8) uint16 {
	c.doDummyMemoryRead(uint16(address))
	return uint16(address + index)
}

func (c *CPU) GetValueByIndirectAbsoluteMode(initialAddress uint16) uint8 {
	lastByte := c.readFromMemory(initialAddress)
	firstByte := c.readFromMemory(initialAddress + 1)

	return c.readFromMemory(uint16(firstByte)<<8 | uint16(lastByte))
}

func (c *CPU) GetValueByIndexedIndirectXMode(initialAddress uint8) uint8 {
	c.readFromMemory(uint16(initialAddress))
	return c.readFromMemory(c.GetAddressByIndexedIndirectXMode(uint16(initialAddress)))
}

func (c *CPU) GetValueByIndirectIndexedYMode(initialAddress uint8) uint8 {
	return c.readFromMemory(c.GetAddressByIndirectIndexedYMode(initialAddress))
}

func (c *CPU) GetAddressByIndexedAbsoluteMode(address uint16, index uint8) uint16 {
	newAddress := address + uint16(index)

	if address&0xFF00 != newAddress&0xFF00 {
		c.doDummyMemoryRead(getProvisionalIndexedAddress(address, newAddress))
	}

	return newAddress
}

func (c *CPU) GetAddressByIndexedAbsoluteModeWithDummyRead(address uint16, index uint8) uint16 {
	newAddress := address + uint16(index)
	c.doDummyMemoryRead(getProvisionalIndexedAddress(address, newAddress))
	return newAddress
}

func (c *CPU) GetAddressByIndexedIndirectXModeWithDummyRead(initialAddress uint16) uint16 {
	c.doDummyMemoryRead(initialAddress)
	return c.GetAddressByIndexedIndirectXMode(initialAddress)
}

func (c *CPU) GetAddressByIndexedIndirectXMode(initialAddress uint16) uint16 {
	pivotAddress := uint8(initialAddress) + c.x
	lastByte := c.readFromMemory(uint16(pivotAddress))
	firstByte := c.readFromMemory(uint16(pivotAddress + 1))

	return uint16(firstByte)<<8 | uint16(lastByte)
}

func (c *CPU) GetAddressByIndirectIndexedYMode(initialAddress uint8) uint16 {
	lowByte := c.readFromMemory(uint16(initialAddress))
	highByte := c.readFromMemory(uint16(initialAddress + 1))
	baseAddress := uint16(highByte)<<8 | uint16(lowByte)
	newAddress := baseAddress + uint16(c.y)

	if baseAddress&0xFF00 != newAddress&0xFF00 {
		c.doDummyMemoryRead(getProvisionalIndexedAddress(baseAddress, newAddress))
	}

	return newAddress
}

func (c *CPU) GetAddressByIndirectIndexedYModeWithDummyRead(initialAddress uint8) uint16 {
	lastByte := c.readFromMemory(uint16(initialAddress))
	firstByte := c.readFromMemory(uint16(initialAddress + 1))
	incompleteAddress := uint16(firstByte)<<8 | uint16(lastByte)
	newAddress := incompleteAddress + uint16(c.y)

	c.doDummyMemoryRead(getProvisionalIndexedAddress(incompleteAddress, newAddress))

	return newAddress
}

func getProvisionalIndexedAddress(baseAddress, effectiveAddress uint16) uint16 {
	return (baseAddress & 0xFF00) | (effectiveAddress & 0x00FF)
}
