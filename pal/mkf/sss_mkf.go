package mkf

import "unsafe"

type SSSMkf struct {
	Mkf
}

func NewSSSMkf(file string) (SSSMkf, error) {
	ret := SSSMkf{Mkf{}}
	return ret, ret.Open(file)
}

func (sm *SSSMkf) GetMsgOffsetChunk() (*MsgOffsetChunk, error) {
	buf, err := sm.ReadChunk(3)
	if err != nil {
		return nil, err
	}
	return &MsgOffsetChunk{NewPlaneChunk(buf)}, nil
}

type EventObject struct {
	VanishTime               SHORT  // vanish time (?)
	X                        WORD   // X coordinate on the map
	Y                        WORD   // Y coordinate on the map
	Layer                    SHORT  // layer value
	TriggerScript            WORD   // Trigger script entry
	AutoScript               WORD   // Auto script entry
	State                    SHORT  // state of this object
	TriggerMode              WORD   // trigger mode
	SpriteNum                WORD   // number of the sprite
	SpriteFrames             USHORT // total number of frames of the sprite
	Direction                WORD   // direction
	CurrentFrameNum          WORD   // current frame number
	ScriptIdleFrame          USHORT // count of idle frames, used by trigger script
	SpritePtrOffset          WORD   // FIXME: ???
	SpriteFramesAuto         USHORT // total number of frames of the sprite, used by auto script
	ScriptIdleFrameCountAuto WORD   // count of idle frames, used by auto script
}

type ScriptEntry struct {
	Operation WORD
	Operand   [3]WORD
}

// chunk 0 in sss.mkf
type EventObjectChunk struct {
	PlaneChunk
}

type MsgOffsetChunk struct {
	PlaneChunk
}

// chunk 4 in sss.mkf
type ScriptEntryChunk struct {
	PlaneChunk
}

func NewScriptEntryChunk(data []byte) ScriptEntryChunk {
	return ScriptEntryChunk{NewPlaneChunk(data)}
}

func (sc *ScriptEntryChunk) GetScriptEntry(idx int) *ScriptEntry {
	return (*ScriptEntry)(sc.Get(idx, unsafe.Sizeof(ScriptEntryChunk{})))
}
