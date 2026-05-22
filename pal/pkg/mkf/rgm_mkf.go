package mkf

// RgmMkf represents the RGM.MKF file which stores character face bitmaps
type RgmMkf struct {
	Mkf
}

// NewRgmMkf creates a new RgmMkf instance
func NewRgmMkf(file string) (RgmMkf, error) {
	ret := RgmMkf{Mkf{}}
	return ret, ret.Open(file)
}

// GetFaceChunk retrieves a character face chunk from RGM.MKF
func (rm *RgmMkf) GetFaceChunk(faceNum INT) (*RgmChunk, error) {
	data, err := rm.ReadChunk(faceNum)
	if err != nil {
		return nil, err
	}
	return NewRgmChunk(data), nil
}

// GetFaceBmp retrieves a character face bitmap from RGM.MKF
func (rm *RgmMkf) GetFaceBmp(faceNum INT) (*BitMap, error) {
	chunk, err := rm.GetFaceChunk(faceNum)
	if err != nil {
		return nil, err
	}
	return chunk.GetBmp(), nil
}

// RgmChunk represents a chunk in RGM.MKF
// RGM.MKF 的 chunk 是直接的 RLE 位图数据，不需要解压
type RgmChunk struct {
	data []byte
}

// NewRgmChunk creates a new RgmChunk from raw data
func NewRgmChunk(data []byte) *RgmChunk {
	return &RgmChunk{data: data}
}

// GetBmp returns the bitmap from the chunk (RGM.MKF 不需要解压)
func (c *RgmChunk) GetBmp() *BitMap {
	if len(c.data) == 0 {
		return nil
	}

	ret := NewRLEBitMap(c.data)
	return ret
}
