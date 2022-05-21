package c64

import _ "embed"

//go:embed roms/kernal
var ROMKernal []byte

//go:embed roms/basic
var ROMBasic []byte

//go:embed roms/chargen
var ROMChargen []byte
