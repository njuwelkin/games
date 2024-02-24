package mkf

const (
	STATUS_BACKGROUND_FBPNUM   = 0
	MAINMENU_BACKGROUND_FBPNUM = 60
)

type FbpMkf struct {
	Mkf
}

func (fm *FbpMkf) GetManMenuBgdChunk() (FbpChunk, error) {
	data, err := fm.ReadChunk(MAINMENU_BACKGROUND_FBPNUM)
	if err != nil {
		return FbpChunk{}, err
	}
	return NewFbpChunk(data), nil
}

func (fm *FbpMkf) GetManMenuBgdBmp() (*BitMap, error) {
	chunk, err := fm.GetManMenuBgdChunk()
	if err != nil {
		return nil, err
	}
	return chunk.GetBmp(), nil
}

type FbpChunk struct{ CompressedChunk }

func NewFbpChunk(data []byte) FbpChunk {
	return FbpChunk{NewCompressedChunk(data)}
}

func (c *FbpChunk) GetBmp() *BitMap {
	data, err := c.Decompress()
	if err != nil {
		return nil
	}
	ret := BitMap{}
	ret.h = 200
	ret.w = 320
	ret.data = data
	return &ret
}
