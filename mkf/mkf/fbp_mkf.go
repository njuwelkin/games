package mkf

const (
	STATUS_BACKGROUND_FBPNUM = 0
)

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
