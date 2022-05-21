package c64

import (
	"crypto/sha512"
	"encoding/hex"
	"testing"
)

func TestInstanciation(t *testing.T) {
	c64 := Make(NTSC)
	c64.Init()
	c64.CPU.Reset()
}

func TestRoms(t *testing.T) {
	const kernalSha512 = "c6ef3021ab08151bd93399ed8c2a97896cb1fb1e2820865622ce1a9169242b48315c71f58aba26b7f720e872b9b941e651378a7a3d99218e1c104d55e412d25c"
	const basicSha512 = "683c3ca9bf14d71b988b35381843a9d8f4e083254b45f2f2a27c1a1a3508090134156171de56604fe3be65ae0d4efbd24e0923d7f6c19f3449c25212711e2320"
	const chargenSha512 = "63b583016b46d5c7af11dd1b57fc721370fbfaf50b7ef2c15a7f40dd420dcbcaf7d3e173d7d42b88e6ed45336bd0cae67b921a15a6b8ee6d0a0925a3ba7caab2"

	kernalHash := sha512.Sum512(ROMKernal)
	if kernalSha512 != hex.EncodeToString(kernalHash[:]) {
		t.Errorf("Kernal ROM hash failed (`roms/kernal`)")
	}
	basicHash := sha512.Sum512(ROMBasic)
	if basicSha512 != hex.EncodeToString(basicHash[:]) {
		t.Errorf("Basic ROM hash failed (`roms/basic`)")
	}
	chargenHash := sha512.Sum512(ROMChargen)
	if chargenSha512 != hex.EncodeToString(chargenHash[:]) {
		t.Errorf("Chargen ROM hash failed (`roms/chargen`)")
	}
}
