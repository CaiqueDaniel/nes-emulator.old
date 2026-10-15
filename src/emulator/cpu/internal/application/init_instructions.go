package application

type instructionSet [256]*instruction

type instruction struct {
	Method    func([]uint8)
	ArgsBytes int
}

func (c *CPU) initInstructions() instructionSet {
	var instructionSet instructionSet

	instructionSet[0xEA] = &instruction{
		Method:    func(u []uint8) { c.NoOp() },
		ArgsBytes: 0,
	}

	instructionSet[0x00] = &instruction{
		Method:    func(u []uint8) { c.Break() },
		ArgsBytes: 1,
	}

	c.appendLDAInstructions(&instructionSet)
	c.appendLDXInstructions(&instructionSet)
	c.appendLDYInstructions(&instructionSet)
	c.appendSTAInstructions(&instructionSet)
	c.appendSTXInstructions(&instructionSet)
	c.appendSTYInstructions(&instructionSet)
	c.appendTransferInstructions(&instructionSet)
	c.appendADCInstructions(&instructionSet)
	c.appendSBCInstructions(&instructionSet)
	c.appendIncrementInstructions(&instructionSet)
	c.appendDecrementInstructions(&instructionSet)

	c.appendShiftInstructions(&instructionSet)
	c.appendBitwiseInstructions(&instructionSet)
	c.appendCompareInstructions(&instructionSet)
	c.appendBranchInstructions(&instructionSet)
	c.appendJumpInstructions(&instructionSet)
	c.appendStackInstructions(&instructionSet)
	c.appendFlagsInstructions(&instructionSet)

	return instructionSet
}

func (c *CPU) appendLDAInstructions(instructionSet *instructionSet) {
	instructionSet[0xA9] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(u[0], ACCUMULATOR) },
		ArgsBytes: 1,
	}
	instructionSet[0xA5] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByZeroPageMode(u[0]), ACCUMULATOR) },
		ArgsBytes: 1,
	}
	instructionSet[0xB5] = &instruction{
		Method: func(u []uint8) {
			c.LoadValueIntoRegister(c.GetValueByZeroPageIndexedModeWithDummyRead(u[0], c.x), ACCUMULATOR)
		},
		ArgsBytes: 1,
	}
	instructionSet[0xAD] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByAbsoluteMode(parseAddress(u)), ACCUMULATOR) },
		ArgsBytes: 2,
	}
	instructionSet[0xBD] = &instruction{
		Method: func(u []uint8) {
			c.LoadValueIntoRegister(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.x), ACCUMULATOR)
		},
		ArgsBytes: 2,
	}
	instructionSet[0xB9] = &instruction{
		Method: func(u []uint8) {
			c.LoadValueIntoRegister(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.y), ACCUMULATOR)
		},
		ArgsBytes: 2,
	}
	instructionSet[0xA1] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByIndexedIndirectXMode(u[0]), ACCUMULATOR) },
		ArgsBytes: 1,
	}
	instructionSet[0xB1] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByIndirectIndexedYMode(u[0]), ACCUMULATOR) },
		ArgsBytes: 1,
	}
}

func (c *CPU) appendLDXInstructions(instructionSet *instructionSet) {
	instructionSet[0xA2] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(u[0], REGISTER_X) },
		ArgsBytes: 1,
	}
	instructionSet[0xA6] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByZeroPageMode(u[0]), REGISTER_X) },
		ArgsBytes: 1,
	}
	instructionSet[0xB6] = &instruction{
		Method: func(u []uint8) {
			c.LoadValueIntoRegister(c.GetValueByZeroPageIndexedModeWithDummyRead(u[0], c.x), REGISTER_X)
		},
		ArgsBytes: 1,
	}
	instructionSet[0xAE] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByAbsoluteMode(parseAddress(u)), REGISTER_X) },
		ArgsBytes: 2,
	}
	instructionSet[0xBE] = &instruction{
		Method: func(u []uint8) {
			c.LoadValueIntoRegister(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.y), REGISTER_X)
		},
		ArgsBytes: 2,
	}
}

func (c *CPU) appendLDYInstructions(instructionSet *instructionSet) {
	instructionSet[0xA0] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(u[0], REGISTER_Y) },
		ArgsBytes: 1,
	}
	instructionSet[0xA4] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByZeroPageMode(u[0]), REGISTER_Y) },
		ArgsBytes: 1,
	}
	instructionSet[0xB4] = &instruction{
		Method: func(u []uint8) {
			c.LoadValueIntoRegister(c.GetValueByZeroPageIndexedModeWithDummyRead(u[0], c.x), REGISTER_Y)
		},
		ArgsBytes: 1,
	}
	instructionSet[0xAC] = &instruction{
		Method:    func(u []uint8) { c.LoadValueIntoRegister(c.GetValueByAbsoluteMode(parseAddress(u)), REGISTER_Y) },
		ArgsBytes: 2,
	}
	instructionSet[0xBC] = &instruction{
		Method: func(u []uint8) {
			c.LoadValueIntoRegister(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.x), REGISTER_Y)
		},
		ArgsBytes: 2,
	}
}

func (c *CPU) appendSTAInstructions(instructionSet *instructionSet) {
	instructionSet[0x85] = &instruction{
		Method:    func(u []uint8) { c.StoreRegisterIntoAbsoluteMemory(uint16(u[0]), ACCUMULATOR) },
		ArgsBytes: 1,
	}

	instructionSet[0x95] = &instruction{
		Method: func(u []uint8) {
			c.StoreRegisterIntoAbsoluteMemory(c.GetAddressByZeroPageIndexedModeWithDummyRead(u[0], c.x), ACCUMULATOR)
		},
		ArgsBytes: 1,
	}

	instructionSet[0x8D] = &instruction{
		Method:    func(u []uint8) { c.StoreRegisterIntoAbsoluteMemory(parseAddress(u), ACCUMULATOR) },
		ArgsBytes: 2,
	}

	instructionSet[0x9D] = &instruction{
		Method: func(u []uint8) {
			c.StoreRegisterIntoAbsoluteMemory(c.GetAddressByIndexedAbsoluteModeWithDummyRead(parseAddress(u), c.x), ACCUMULATOR)
		},
		ArgsBytes: 2,
	}

	instructionSet[0x99] = &instruction{
		Method: func(u []uint8) {
			c.StoreRegisterIntoAbsoluteMemory(c.GetAddressByIndexedAbsoluteModeWithDummyRead(parseAddress(u), c.y), ACCUMULATOR)
		},
		ArgsBytes: 2,
	}

	instructionSet[0x81] = &instruction{
		Method: func(u []uint8) {
			c.StoreRegisterIntoAbsoluteMemory(c.GetAddressByIndexedIndirectXModeWithDummyRead(uint16(u[0])), ACCUMULATOR)
		},
		ArgsBytes: 1,
	}

	instructionSet[0x91] = &instruction{
		Method: func(u []uint8) {
			c.StoreRegisterIntoAbsoluteMemory(c.GetAddressByIndirectIndexedYModeWithDummyRead(u[0]), ACCUMULATOR)
		},
		ArgsBytes: 1,
	}

}

func (c *CPU) appendSTXInstructions(instructionSet *instructionSet) {
	instructionSet[0x86] = &instruction{
		Method:    func(u []uint8) { c.StoreRegisterIntoAbsoluteMemory(uint16(u[0]), REGISTER_X) },
		ArgsBytes: 1,
	}

	instructionSet[0x96] = &instruction{
		Method: func(u []uint8) {
			c.StoreRegisterIntoAbsoluteMemory(c.GetAddressByZeroPageIndexedModeWithDummyRead(u[0], c.y), REGISTER_X)
		},
		ArgsBytes: 1,
	}

	instructionSet[0x8E] = &instruction{
		Method:    func(u []uint8) { c.StoreRegisterIntoAbsoluteMemory(parseAddress(u), REGISTER_X) },
		ArgsBytes: 2,
	}
}

func (c *CPU) appendSTYInstructions(instructionSet *instructionSet) {
	instructionSet[0x84] = &instruction{
		Method:    func(u []uint8) { c.StoreRegisterIntoAbsoluteMemory(uint16(u[0]), REGISTER_Y) },
		ArgsBytes: 1,
	}

	instructionSet[0x94] = &instruction{
		Method: func(u []uint8) {
			c.StoreRegisterIntoAbsoluteMemory(c.GetAddressByZeroPageIndexedModeWithDummyRead(u[0], c.x), REGISTER_Y)
		},
		ArgsBytes: 1,
	}

	instructionSet[0x8C] = &instruction{
		Method:    func(u []uint8) { c.StoreRegisterIntoAbsoluteMemory(parseAddress(u), REGISTER_Y) },
		ArgsBytes: 2,
	}
}

func (c *CPU) appendTransferInstructions(instructionSet *instructionSet) {
	instructionSet[0xAA] = &instruction{
		Method:    func(u []uint8) { c.TransferFromAccumulatorToRegister(REGISTER_X) },
		ArgsBytes: 0,
	}

	instructionSet[0xA8] = &instruction{
		Method:    func(u []uint8) { c.TransferFromAccumulatorToRegister(REGISTER_Y) },
		ArgsBytes: 0,
	}

	instructionSet[0xBA] = &instruction{
		Method:    func(u []uint8) { c.TransferStackToX() },
		ArgsBytes: 0,
	}

	instructionSet[0x8A] = &instruction{
		Method:    func(u []uint8) { c.TransferFromRegisterToAccumulator(REGISTER_X) },
		ArgsBytes: 0,
	}

	instructionSet[0x98] = &instruction{
		Method:    func(u []uint8) { c.TransferFromRegisterToAccumulator(REGISTER_Y) },
		ArgsBytes: 0,
	}

	instructionSet[0x9A] = &instruction{
		Method:    func(u []uint8) { c.TransferXToStack() },
		ArgsBytes: 0,
	}
}

func (c *CPU) appendADCInstructions(instructionSet *instructionSet) {
	instructionSet[0x69] = &instruction{
		Method:    func(u []uint8) { c.AddWithCarry(u[0]) },
		ArgsBytes: 1,
	}

	instructionSet[0x65] = &instruction{
		Method:    func(u []uint8) { c.AddWithCarry(c.GetValueByZeroPageMode(u[0])) },
		ArgsBytes: 1,
	}

	instructionSet[0x75] = &instruction{
		Method:    func(u []uint8) { c.AddWithCarry(c.GetValueByZeroPageIndexedModeWithDummyRead(u[0], c.x)) },
		ArgsBytes: 1,
	}

	instructionSet[0x6D] = &instruction{
		Method:    func(u []uint8) { c.AddWithCarry(c.GetValueByAbsoluteMode(parseAddress(u))) },
		ArgsBytes: 2,
	}

	instructionSet[0x7D] = &instruction{
		Method:    func(u []uint8) { c.AddWithCarry(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.x)) },
		ArgsBytes: 2,
	}

	instructionSet[0x79] = &instruction{
		Method:    func(u []uint8) { c.AddWithCarry(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.y)) },
		ArgsBytes: 2,
	}

	instructionSet[0x61] = &instruction{
		Method:    func(u []uint8) { c.AddWithCarry(c.GetValueByIndexedIndirectXMode(u[0])) },
		ArgsBytes: 1,
	}

	instructionSet[0x71] = &instruction{
		Method:    func(u []uint8) { c.AddWithCarry(c.GetValueByIndirectIndexedYMode(u[0])) },
		ArgsBytes: 1,
	}

}

func (c *CPU) appendSBCInstructions(instructionSet *instructionSet) {
	instructionSet[0xE9] = &instruction{
		Method:    func(u []uint8) { c.SubtractWithCarry(u[0]) },
		ArgsBytes: 1,
	}

	instructionSet[0xE5] = &instruction{
		Method:    func(u []uint8) { c.SubtractWithCarry(c.GetValueByZeroPageMode(u[0])) },
		ArgsBytes: 1,
	}

	instructionSet[0xF5] = &instruction{
		Method:    func(u []uint8) { c.SubtractWithCarry(c.GetValueByZeroPageIndexedModeWithDummyRead(u[0], c.x)) },
		ArgsBytes: 1,
	}

	instructionSet[0xED] = &instruction{
		Method:    func(u []uint8) { c.SubtractWithCarry(c.GetValueByAbsoluteMode(parseAddress(u))) },
		ArgsBytes: 2,
	}

	instructionSet[0xFD] = &instruction{
		Method:    func(u []uint8) { c.SubtractWithCarry(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.x)) },
		ArgsBytes: 2,
	}

	instructionSet[0xF9] = &instruction{
		Method:    func(u []uint8) { c.SubtractWithCarry(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.y)) },
		ArgsBytes: 2,
	}

	instructionSet[0xE1] = &instruction{
		Method:    func(u []uint8) { c.SubtractWithCarry(c.GetValueByIndexedIndirectXMode(u[0])) },
		ArgsBytes: 1,
	}

	instructionSet[0xF1] = &instruction{
		Method:    func(u []uint8) { c.SubtractWithCarry(c.GetValueByIndirectIndexedYMode(u[0])) },
		ArgsBytes: 1,
	}

}

func (c *CPU) appendIncrementInstructions(instructionSet *instructionSet) {
	instructionSet[0xE6] = &instruction{
		Method:    func(u []uint8) { c.IncrementMemory(uint16(u[0])) },
		ArgsBytes: 1,
	}

	instructionSet[0xF6] = &instruction{
		Method:    func(u []uint8) { c.IncrementMemory(c.GetAddressByZeroPageIndexedModeWithDummyRead(u[0], c.x)) },
		ArgsBytes: 1,
	}

	instructionSet[0xEE] = &instruction{
		Method:    func(u []uint8) { c.IncrementMemory(parseAddress(u)) },
		ArgsBytes: 2,
	}

	instructionSet[0xFE] = &instruction{
		Method: func(u []uint8) {
			c.IncrementMemory(c.GetAddressByIndexedAbsoluteModeWithDummyRead(parseAddress(u), c.x))
		},
		ArgsBytes: 2,
	}

	instructionSet[0xE8] = &instruction{
		Method:    func(u []uint8) { c.IncrementRegister(REGISTER_X) },
		ArgsBytes: 0,
	}

	instructionSet[0xC8] = &instruction{
		Method:    func(u []uint8) { c.IncrementRegister(REGISTER_Y) },
		ArgsBytes: 0,
	}
}

func (c *CPU) appendDecrementInstructions(instructionSet *instructionSet) {
	instructionSet[0xC6] = &instruction{
		Method:    func(u []uint8) { c.DecrementMemory(uint16(u[0])) },
		ArgsBytes: 1,
	}

	instructionSet[0xD6] = &instruction{
		Method:    func(u []uint8) { c.DecrementMemory(c.GetAddressByZeroPageIndexedModeWithDummyRead(u[0], c.x)) },
		ArgsBytes: 1,
	}

	instructionSet[0xCE] = &instruction{
		Method:    func(u []uint8) { c.DecrementMemory(parseAddress(u)) },
		ArgsBytes: 2,
	}

	instructionSet[0xDE] = &instruction{
		Method: func(u []uint8) {
			c.DecrementMemory(c.GetAddressByIndexedAbsoluteModeWithDummyRead(parseAddress(u), c.x))
		},
		ArgsBytes: 2,
	}

	instructionSet[0xCA] = &instruction{
		Method:    func(u []uint8) { c.DecrementRegister(REGISTER_X) },
		ArgsBytes: 0,
	}

	instructionSet[0x88] = &instruction{
		Method:    func(u []uint8) { c.DecrementRegister(REGISTER_Y) },
		ArgsBytes: 0,
	}
}

func (c *CPU) appendShiftInstructions(instructionSet *instructionSet) {
	// ASL
	instructionSet[0x0A] = &instruction{Method: func(u []uint8) { c.ArithmeticShiftLeft() }, ArgsBytes: 0}
	instructionSet[0x06] = &instruction{Method: func(u []uint8) { c.ArithmeticShiftLeftZeroPage(u[0], false) }, ArgsBytes: 1}
	instructionSet[0x16] = &instruction{Method: func(u []uint8) { c.ArithmeticShiftLeftZeroPage(u[0], true) }, ArgsBytes: 1}
	instructionSet[0x0E] = &instruction{Method: func(u []uint8) { c.ArithmeticShiftLeftAbsolute(parseAddress(u), false) }, ArgsBytes: 2}
	instructionSet[0x1E] = &instruction{Method: func(u []uint8) { c.ArithmeticShiftLeftAbsolute(parseAddress(u), true) }, ArgsBytes: 2}

	// LSR
	instructionSet[0x4A] = &instruction{Method: func(u []uint8) { c.LogicalShiftRight() }, ArgsBytes: 0}
	instructionSet[0x46] = &instruction{Method: func(u []uint8) { c.LogicalShiftRightZeroPage(u[0], false) }, ArgsBytes: 1}
	instructionSet[0x56] = &instruction{Method: func(u []uint8) { c.LogicalShiftRightZeroPage(u[0], true) }, ArgsBytes: 1}
	instructionSet[0x4E] = &instruction{Method: func(u []uint8) { c.LogicalShiftRightAbsolute(parseAddress(u), false) }, ArgsBytes: 2}
	instructionSet[0x5E] = &instruction{Method: func(u []uint8) { c.LogicalShiftRightAbsolute(parseAddress(u), true) }, ArgsBytes: 2}

	// ROL
	instructionSet[0x2A] = &instruction{Method: func(u []uint8) { c.RotateLeft() }, ArgsBytes: 0}
	instructionSet[0x26] = &instruction{Method: func(u []uint8) { c.RotateLeftZeroPage(u[0], false) }, ArgsBytes: 1}
	instructionSet[0x36] = &instruction{Method: func(u []uint8) { c.RotateLeftZeroPage(u[0], true) }, ArgsBytes: 1}
	instructionSet[0x2E] = &instruction{Method: func(u []uint8) { c.RotateLeftAbsolute(parseAddress(u), false) }, ArgsBytes: 2}
	instructionSet[0x3E] = &instruction{Method: func(u []uint8) { c.RotateLeftAbsolute(parseAddress(u), true) }, ArgsBytes: 2}

	// ROR
	instructionSet[0x6A] = &instruction{Method: func(u []uint8) { c.RotateRight() }, ArgsBytes: 0}
	instructionSet[0x66] = &instruction{Method: func(u []uint8) { c.RotateRightZeroPage(u[0], false) }, ArgsBytes: 1}
	instructionSet[0x76] = &instruction{Method: func(u []uint8) { c.RotateRightZeroPage(u[0], true) }, ArgsBytes: 1}
	instructionSet[0x6E] = &instruction{Method: func(u []uint8) { c.RotateRightAbsolute(parseAddress(u), false) }, ArgsBytes: 2}
	instructionSet[0x7E] = &instruction{Method: func(u []uint8) { c.RotateRightAbsolute(parseAddress(u), true) }, ArgsBytes: 2}
}

func (c *CPU) appendBitwiseInstructions(instructionSet *instructionSet) {
	// AND
	instructionSet[0x29] = &instruction{Method: func(u []uint8) { c.And(u[0]) }, ArgsBytes: 1}
	instructionSet[0x25] = &instruction{Method: func(u []uint8) { c.And(c.GetValueByZeroPageMode(u[0])) }, ArgsBytes: 1}
	instructionSet[0x35] = &instruction{Method: func(u []uint8) { c.And(c.GetValueByZeroPageIndexedModeWithDummyRead(u[0], c.x)) }, ArgsBytes: 1}
	instructionSet[0x2D] = &instruction{Method: func(u []uint8) { c.And(c.GetValueByAbsoluteMode(parseAddress(u))) }, ArgsBytes: 2}
	instructionSet[0x3D] = &instruction{Method: func(u []uint8) { c.And(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.x)) }, ArgsBytes: 2}
	instructionSet[0x39] = &instruction{Method: func(u []uint8) { c.And(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.y)) }, ArgsBytes: 2}
	instructionSet[0x21] = &instruction{Method: func(u []uint8) { c.And(c.GetValueByIndexedIndirectXMode(u[0])) }, ArgsBytes: 1}
	instructionSet[0x31] = &instruction{Method: func(u []uint8) { c.And(c.GetValueByIndirectIndexedYMode(u[0])) }, ArgsBytes: 1}

	// ORA
	instructionSet[0x09] = &instruction{Method: func(u []uint8) { c.Or(u[0]) }, ArgsBytes: 1}
	instructionSet[0x05] = &instruction{Method: func(u []uint8) { c.Or(c.GetValueByZeroPageMode(u[0])) }, ArgsBytes: 1}
	instructionSet[0x15] = &instruction{Method: func(u []uint8) { c.Or(c.GetValueByZeroPageIndexedModeWithDummyRead(u[0], c.x)) }, ArgsBytes: 1}
	instructionSet[0x0D] = &instruction{Method: func(u []uint8) { c.Or(c.GetValueByAbsoluteMode(parseAddress(u))) }, ArgsBytes: 2}
	instructionSet[0x1D] = &instruction{Method: func(u []uint8) { c.Or(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.x)) }, ArgsBytes: 2}
	instructionSet[0x19] = &instruction{Method: func(u []uint8) { c.Or(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.y)) }, ArgsBytes: 2}
	instructionSet[0x01] = &instruction{Method: func(u []uint8) { c.Or(c.GetValueByIndexedIndirectXMode(u[0])) }, ArgsBytes: 1}
	instructionSet[0x11] = &instruction{Method: func(u []uint8) { c.Or(c.GetValueByIndirectIndexedYMode(u[0])) }, ArgsBytes: 1}

	// EOR
	instructionSet[0x49] = &instruction{Method: func(u []uint8) { c.Xor(u[0]) }, ArgsBytes: 1}
	instructionSet[0x45] = &instruction{Method: func(u []uint8) { c.Xor(c.GetValueByZeroPageMode(u[0])) }, ArgsBytes: 1}
	instructionSet[0x55] = &instruction{Method: func(u []uint8) { c.Xor(c.GetValueByZeroPageIndexedModeWithDummyRead(u[0], c.x)) }, ArgsBytes: 1}
	instructionSet[0x4D] = &instruction{Method: func(u []uint8) { c.Xor(c.GetValueByAbsoluteMode(parseAddress(u))) }, ArgsBytes: 2}
	instructionSet[0x5D] = &instruction{Method: func(u []uint8) { c.Xor(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.x)) }, ArgsBytes: 2}
	instructionSet[0x59] = &instruction{Method: func(u []uint8) { c.Xor(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.y)) }, ArgsBytes: 2}
	instructionSet[0x41] = &instruction{Method: func(u []uint8) { c.Xor(c.GetValueByIndexedIndirectXMode(u[0])) }, ArgsBytes: 1}
	instructionSet[0x51] = &instruction{Method: func(u []uint8) { c.Xor(c.GetValueByIndirectIndexedYMode(u[0])) }, ArgsBytes: 1}

	// BIT
	instructionSet[0x24] = &instruction{Method: func(u []uint8) { c.Bit(c.GetValueByZeroPageMode(u[0])) }, ArgsBytes: 1}
	instructionSet[0x2C] = &instruction{Method: func(u []uint8) { c.Bit(c.GetValueByAbsoluteMode(parseAddress(u))) }, ArgsBytes: 2}
}

func (c *CPU) appendCompareInstructions(instructionSet *instructionSet) {
	// CMP
	instructionSet[0xC9] = &instruction{Method: func(u []uint8) { c.CompareWithRegister(u[0], ACCUMULATOR) }, ArgsBytes: 1}
	instructionSet[0xC5] = &instruction{Method: func(u []uint8) { c.CompareWithRegister(c.GetValueByZeroPageMode(u[0]), ACCUMULATOR) }, ArgsBytes: 1}
	instructionSet[0xD5] = &instruction{Method: func(u []uint8) {
		c.CompareWithRegister(c.GetValueByZeroPageIndexedModeWithDummyRead(u[0], c.x), ACCUMULATOR)
	}, ArgsBytes: 1}
	instructionSet[0xCD] = &instruction{Method: func(u []uint8) { c.CompareWithRegister(c.GetValueByAbsoluteMode(parseAddress(u)), ACCUMULATOR) }, ArgsBytes: 2}
	instructionSet[0xDD] = &instruction{Method: func(u []uint8) {
		c.CompareWithRegister(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.x), ACCUMULATOR)
	}, ArgsBytes: 2}
	instructionSet[0xD9] = &instruction{Method: func(u []uint8) {
		c.CompareWithRegister(c.GetValueByIndexedAbsoluteMode(parseAddress(u), c.y), ACCUMULATOR)
	}, ArgsBytes: 2}
	instructionSet[0xC1] = &instruction{Method: func(u []uint8) { c.CompareWithRegister(c.GetValueByIndexedIndirectXMode(u[0]), ACCUMULATOR) }, ArgsBytes: 1}
	instructionSet[0xD1] = &instruction{Method: func(u []uint8) { c.CompareWithRegister(c.GetValueByIndirectIndexedYMode(u[0]), ACCUMULATOR) }, ArgsBytes: 1}

	// CPX
	instructionSet[0xE0] = &instruction{Method: func(u []uint8) { c.CompareWithRegister(u[0], REGISTER_X) }, ArgsBytes: 1}
	instructionSet[0xE4] = &instruction{Method: func(u []uint8) { c.CompareWithRegister(c.GetValueByZeroPageMode(u[0]), REGISTER_X) }, ArgsBytes: 1}
	instructionSet[0xEC] = &instruction{Method: func(u []uint8) { c.CompareWithRegister(c.GetValueByAbsoluteMode(parseAddress(u)), REGISTER_X) }, ArgsBytes: 2}

	// CPY
	instructionSet[0xC0] = &instruction{Method: func(u []uint8) { c.CompareWithRegister(u[0], REGISTER_Y) }, ArgsBytes: 1}
	instructionSet[0xC4] = &instruction{Method: func(u []uint8) { c.CompareWithRegister(c.GetValueByZeroPageMode(u[0]), REGISTER_Y) }, ArgsBytes: 1}
	instructionSet[0xCC] = &instruction{Method: func(u []uint8) { c.CompareWithRegister(c.GetValueByAbsoluteMode(parseAddress(u)), REGISTER_Y) }, ArgsBytes: 2}
}

func (c *CPU) appendBranchInstructions(instructionSet *instructionSet) {
	instructionSet[0x90] = &instruction{Method: func(u []uint8) { c.BranchIfCarryIsClear(u[0]) }, ArgsBytes: 1}
	instructionSet[0xB0] = &instruction{Method: func(u []uint8) { c.BranchIfCarryIsSet(u[0]) }, ArgsBytes: 1}
	instructionSet[0xF0] = &instruction{Method: func(u []uint8) { c.BranchIfEqual(u[0]) }, ArgsBytes: 1}
	instructionSet[0xD0] = &instruction{Method: func(u []uint8) { c.BranchIfNotEqual(u[0]) }, ArgsBytes: 1}
	instructionSet[0x30] = &instruction{Method: func(u []uint8) { c.BranchIfNegative(u[0]) }, ArgsBytes: 1}
	instructionSet[0x10] = &instruction{Method: func(u []uint8) { c.BranchIfPositive(u[0]) }, ArgsBytes: 1}
	instructionSet[0x50] = &instruction{Method: func(u []uint8) { c.BranchIfOverflowClear(u[0]) }, ArgsBytes: 1}
	instructionSet[0x70] = &instruction{Method: func(u []uint8) { c.BranchIfOverflowSet(u[0]) }, ArgsBytes: 1}
}

func (c *CPU) appendJumpInstructions(instructionSet *instructionSet) {
	instructionSet[0x4C] = &instruction{Method: func(u []uint8) { c.JumpProgramCounterToValue(parseAddress(u)) }, ArgsBytes: 2}
	instructionSet[0x6C] = &instruction{Method: func(u []uint8) { c.JumpProgramCounterByIndirectValue(parseAddress(u)) }, ArgsBytes: 2}
	instructionSet[0x20] = &instruction{Method: func(u []uint8) { c.JumpProgramCounterToSubRoutine(parseAddress(u)) }, ArgsBytes: 2}
	instructionSet[0x60] = &instruction{Method: func(u []uint8) { c.ReturnFromSubRoutine() }, ArgsBytes: 0}
	instructionSet[0x40] = &instruction{Method: func(u []uint8) { c.ReturnFromInterrupt() }, ArgsBytes: 0}
}

func (c *CPU) appendStackInstructions(instructionSet *instructionSet) {
	instructionSet[0x48] = &instruction{Method: func(u []uint8) { c.PushAccToStack() }, ArgsBytes: 0}
	instructionSet[0x68] = &instruction{Method: func(u []uint8) { c.PullAccFromStack() }, ArgsBytes: 0}
	instructionSet[0x08] = &instruction{Method: func(u []uint8) { c.PushStatusIntoStack() }, ArgsBytes: 0}
	instructionSet[0x28] = &instruction{Method: func(u []uint8) { c.PullStatusFromStack() }, ArgsBytes: 0}
}

func (c *CPU) appendFlagsInstructions(instructionSet *instructionSet) {
	instructionSet[0x18] = &instruction{Method: func(u []uint8) { c.ClearCarryFlag() }, ArgsBytes: 0}
	instructionSet[0x38] = &instruction{Method: func(u []uint8) { c.SetCarryFlag() }, ArgsBytes: 0}
	instructionSet[0x58] = &instruction{Method: func(u []uint8) { c.ClearInterruptFlag() }, ArgsBytes: 0}
	instructionSet[0x78] = &instruction{Method: func(u []uint8) { c.SetInterruptFlag() }, ArgsBytes: 0}
	instructionSet[0xD8] = &instruction{Method: func(u []uint8) { c.ClearDecimalFlag() }, ArgsBytes: 0}
	instructionSet[0xF8] = &instruction{Method: func(u []uint8) { c.SetDecimalFlag() }, ArgsBytes: 0}
	instructionSet[0xB8] = &instruction{Method: func(u []uint8) { c.ClearOverflowFlag() }, ArgsBytes: 0}
}

func parseAddress(bytes []uint8) uint16 {
	if len(bytes) != 2 {
		return 0
	}

	return (uint16(bytes[0]) << 8) + uint16(bytes[1])
}
