package c64

import (
	"errors"
	"os"
	"strconv"
)

// Reads a 16 bit int from 2 bytes little endian wise
func getUint16(from []byte) uint16 {
	return (uint16(from[1]) << 8) | uint16(from[0])
}

// Reads a PRG file from a Host OS file and loads it into memory at address specified by the first
// 2 bytes of the file, little endian.
// Returns the address where it was loaded so it can be JMP'ed to.
// If load address is 0x801, it's assumed to be a BASIC program and thus a SYS command is searched for and its address
// is returned instead. (WIP: this is temporary to avoid running the BASIC program)
// WIP Not tested (manual test passed)
func (c64 *C64) LoadPRG(file *os.File) (uint16, error) {
	finfo, err := file.Stat()
	if err != nil {
		return 0, err
	}
	header := make([]byte, 2)
	count, err := file.Read(header)
	if err != nil {
		return 0, err
	} else if count != 2 {
		return 0, errors.New("couldn't read 2-byte load address from file header")
	}
	address := getUint16(header)
	prgSize := uint16(finfo.Size() - 2)
	// If the IO registers bank is visible, make sure that loading the file won't overflow into the IO area because that
	// would produce undefined behavior. (WIP maybe some hack uses this?)
	if c64.isIOon() &&
		((address >= 0xD000 && address <= 0xDFFF) ||
			(address < 0xD000 && address + prgSize > 0xD000)) {
		return 0, errors.New("loading file would overflow into IO area")
	}
	// Read the PRG file minus its header into RAM
	file.ReadAt(c64.RAM[address : address + prgSize], 2)

	// If load address is not BASIC program area, file is assumed to be pure machine code, just JMP to address
	basicArea := c64.ReadUint16(0x2B) // Pointer to start of BASIC area, default is $801
	if address != basicArea {
		return address, nil
	}

	ptr := c64.RAM[address+4:]
	// Else scan BASIC program for SYS command to get address of the start of the machine code program
	// WIP For now we just assume a loader that does a SYS in the first line
	if ptr[0] != 0x9E { // $9E token is Basic SYS command
		return 0, errors.New("unrecognized BASIC loader")
	}

	strAddress := ""
	for ptr[0] != 0 { // Character \0 marks the end of the BASIC line
		ptr = ptr[1:]
		if ptr[0] >= '0' && ptr[0] <= '9' {
			strAddress += string(ptr[0])
		}
	}
	jmpTo, _ := strconv.Atoi(strAddress)
	return uint16(jmpTo), nil
}
