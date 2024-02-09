package mkf

import "unsafe"

type ScriptEntry struct {
	Operation uint16
	Operand   [3]uint16
}

type SssChunk struct {
	PlaneChunk
}

func NewSssChunk(data []byte) SssChunk {
	return SssChunk{NewPlaneChunk(data)}
}

func (sc *SssChunk) GetEntry(idx int) *ScriptEntry {
	return (*ScriptEntry)(sc.Get(idx, unsafe.Sizeof(SssChunk{})))
}
