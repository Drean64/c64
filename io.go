package c64

/**
 *	<io.go> Handles all writes and reads in the I/O memory mapped registers $D000..$DFFF
 *	Mirrored memory, like $D000..$D03F which is mirrored up to $D3FF, is only ever read and written from the original,
 *	unmirrored addresses, I.E. $D040 will never be read from or written  to through these functions, trying to do so
 *	will actually read from or write to $D000.
 *	Reading and writing to those mirror addresses is up to mods and has to be done through direct access to C64.IO[]
 */

/** TODO
 *	Fijandome en VICE, si lees $D020, los 4 high bits vienen todos en 1.
 *	Sospecho que todos los bits no usados siempre se leen como 1, tener en cuenta.
 *	Como las direcciones tipo $D02F que siempre devuelven $FF.
 */

// Read I/O ports
// address is assumed to be $D000..$DFFF
func (c64 *C64) ReadIO(address uint16) byte {
	if address <= 0xD3FF { // VIC video registers
		// WIP test pending
		address &= 0x3F // de-mirror address, registers $D000..$D03F are repeated up to $D3FF
		if address >= 0x2F {
			return 0xFF //  $D02F..$D03F Unused bytes, always read as $FF
		} else if address >= 0x20 {
			// Registers $D020..$D02E are colors and only use low 4 bits. High bits always read as 1
			return c64.IO[address] | 0b11110000
		}
		return c64.IO[address]
	} else {
		// Generic IO read WIP
		// Transpose address into space $0..$FFF of IO bank
		return c64.IO[address & 0xFFF]
	}
}

// Write I/O ports
// address is assumed to be $D000..$DFFF
func (c64 *C64) WriteIO(address uint16, value byte) {
	if address <= 0xD3FF { // VIC video registers
		// WIP test pending
		address &= 0x3F     // de-mirror address, registers $D000..$D03F are repeated up to $D3FF
		if address < 0x20 { // Various VIC registers
			c64.IO[address] = value
		} else if address < 0x2F { // 4-bit color registers
			c64.IO[address] = 0b11110000 | value // higher 4 bits are always 1
		}
		// else unusable bytes in $D02F..$D03F,  ignore the write.
	} else { // Generic IO write WIP
		// Transpose address into space $0..$FFF of IO bank
		c64.IO[address & 0xFFF] = value
		if address == 0xDD00 {
			c64.Vic.setBank(^(value & 0b11)) // Bitwise not of bits 0 & 1 turns 0,1,2,3 into 3,2,1,0
		}
	}
}
