package mkf

const (
	SPRITENUM_SPLASH_TITLE = 0x47
	SPRITENUM_SPLASH_CRANE = 0x49
)

type MgoMkf struct {
	Mkf
}

func NewMgoMkf(file string) (MgoMkf, error) {
	ret := MgoMkf{Mkf{}}
	return ret, ret.Open(file)
}

func (mm *MgoMkf) GetChunk(chunkNum INT) (MgoChunk, error) {
	buf, err := mm.ReadChunk(chunkNum)
	if err != nil {
		return MgoChunk{}, err
	}
	cc := NewCompressedChunk(buf)
	buf, err = cc.Decompress()
	if err != nil {
		return MgoChunk{}, err
	}
	fc := NewBitMapChunk(buf)
	fc.SetFrameSize(32000)
	return MgoChunk{fc}, nil
}

type MgoChunk struct{ BitMapChunk }
