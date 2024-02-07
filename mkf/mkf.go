package main

import (
	"encoding/binary"
	"fmt"
	"image/color"
	"image/color/palette"
	"io"
	"os"
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2"
)

type INT uint32

/*
	type mfkData struct {
		count	uint32
		offset	[count+1]uint32
		data	[count] struct {
			...
		}
	}
*/
type Mkf struct {
	file *os.File
}

func (mkf *Mkf) Open(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	mkf.file = file
	return mkf.init()
}

func (mkf *Mkf) Close() {
	mkf.file.Close()
}

func (mkf *Mkf) init() error {

	return nil
}

func (mkf *Mkf) GetChunkCount() (INT, error) {
	res, err := mkf.readINT(0)
	return (res - 4) >> 2, err
}

func (mkf *Mkf) GetChunkOffset(chunkNum INT) (INT, error) {
	count, err := mkf.GetChunkCount()
	if err != nil {
		return 0, err
	}
	if chunkNum > count {
		return 0, fmt.Errorf("")
	}
	return mkf.readINT(4 * chunkNum)
}

func (mkf *Mkf) GetChunkSize(chunkNum INT) (INT, error) {
	chunkOffset, err := mkf.GetChunkOffset(chunkNum)
	if err != nil {
		return 0, err
	}
	nextChunkOffset, err := mkf.GetChunkOffset(chunkNum + 1)
	if err != nil {
		return 0, err
	}
	//fmt.Println(chunkOffset, nextChunkOffset)
	return nextChunkOffset - chunkOffset, nil
}

func (mkf *Mkf) ReadChunk(chunkNum INT) ([]byte, error) {
	offset, err := mkf.GetChunkOffset(chunkNum)
	if err != nil {
		return nil, err
	}
	nextOffset, err := mkf.GetChunkOffset(chunkNum + 1)
	if err != nil {
		return nil, err
	}

	ret := make([]byte, nextOffset-offset)
	err = mkf.read(offset, ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (mkf *Mkf) read(offset INT, buf []byte) error {
	f := mkf.file
	size := len(buf)
	_, err := f.Seek(int64(offset), io.SeekStart)
	if err != nil {
		return err
	}
	n, err := f.Read(buf)
	if err != nil {
		return err
	}
	if n != int(size) {
		return fmt.Errorf("")
	}
	return nil
}

func (mkf *Mkf) readINT(offset INT) (INT, error) {
	size := unsafe.Sizeof(INT(0))
	buf := make([]byte, size)

	err := mkf.read(offset, buf)
	if err != nil {
		return 0, err
	}
	return INT(binary.LittleEndian.Uint32(buf)), nil
}

type FrameChunk struct {
	data []byte
}

func (c *FrameChunk) GetCount() INT {
	return INT(c.data[0]) | INT(c.data[1])<<8
}

func (c *FrameChunk) getOffset(frameNum INT) INT {
	tile := c.data
	frameNum <<= 1
	return ((INT(tile[frameNum]) | (INT(tile[frameNum+1]) << 8)) << 1)
}

func (c *FrameChunk) GetFrame(frameNum INT) ([]byte, error) {
	count := c.GetCount()
	fmt.Printf("%d images in total\n", count)
	if frameNum >= count {
		return nil, fmt.Errorf("")
	}
	data := c.data
	offset := c.getOffset(frameNum)
	nextOffset := c.getOffset(frameNum + 1)
	fmt.Printf("%d, %d\n", offset, nextOffset)
	//if offset == 0x18444 {
	//	offset = WORD(offset)
	//}
	if nextOffset == 0 {
		return nil, fmt.Errorf("")
	}
	return data[offset:nextOffset], nil
}

/*
type tileBitMapChunkData struct {
	count	uint16
	offset	[count] uint16
	data	[count] rleBitMap
}
*/

type BitMapChunk struct {
	FrameChunk
}

// PAL_SpriteGetFrame
func (bc *BitMapChunk) GetTileBitMap(frameNum INT) (*BitMap, error) {
	frame, err := bc.GetFrame(frameNum)
	if err != nil {
		return nil, err
	}
	return &BitMap{frame}, nil
}

/*
type RLEBitMapData struct {
	width	uint16
	height	uint16
	data ...
}
*/

type BitMap struct {
	data []byte
}

func (bmp *BitMap) GetWidth() INT {
	offset := 0
	data := bmp.data
	if data[0] == 0x02 && data[1] == 0x00 &&
		data[2] == 0x00 && data[3] == 0x00 {
		offset += 4
	}
	return INT(data[offset]) | INT(data[offset+1])<<8
}

func (bmp *BitMap) GetHeight() INT {
	offset := 0
	data := bmp.data
	if data[0] == 0x02 && data[1] == 0x00 &&
		data[2] == 0x00 && data[3] == 0x00 {
		offset += 4
	}
	return INT(data[offset+2]) | INT(data[offset+3])<<8
}

func (bmp *BitMap) GetNINT(n INT) INT {
	offset := 0
	data := bmp.data
	if data[0] == 0x02 && data[1] == 0x00 &&
		data[2] == 0x00 && data[3] == 0x00 {
		offset += 4
	}
	return INT(data[offset+int(2*n)]) | INT(data[offset+int(2*n)+1])<<8
}

func (bmp *BitMap) PrintRaw() {
	for i := 5; i < len(bmp.data); i++ {
		fmt.Printf("%x ", bmp.data[i])
	}
	fmt.Println()
}

func (bmp *BitMap) Decode() []byte {
	//w := bmp.GetWidth()
	//h := bmp.GetHeight()
	//l := w * h
	//var uiSrcX INT = 0

	data := bmp.data[5:]
	tIdx := 0
	var countPix INT = 0
	var countEmpty INT = 0
	for tIdx < len(data) && data[tIdx] != 0 {
		T := INT(data[tIdx])

		if T&0x80 != 0 { //&& T <= 0x80+w {
			countEmpty += T - 0x80
			tIdx++
		} else {
			countPix += T
			tIdx += int(T)
			tIdx++
		}
	}
	//fmt.Println(countPix, countEmpty)
	return nil
}

func (bmp *BitMap) ToImage() *ebiten.Image {
	return bmp.toImage(palette.Plan9)
}

func (bmp *BitMap) toImage(plt []color.Color) *ebiten.Image {
	w := int(bmp.GetWidth())
	h := int(bmp.GetHeight())
	//l := w * h
	img := ebiten.NewImage(int(w), int(h))

	//var uiSrcX INT = 0

	data := bmp.data[5:]
	tIdx := 0
	x, y := 0, 0
	for tIdx < len(data) && data[tIdx] != 0 {
		T := INT(data[tIdx])

		if T&0x80 != 0 { //&& T <= 0x80+w {
			x += int(T - 0x80)
			y += x / w
			x %= w

			tIdx++
		} else {
			tIdx++
			if tIdx >= len(data) {
				break
			}
			for j := 0; j < int(T); j++ {
				img.Set(x, y, PixToRGBA(data[tIdx], plt))
				x, y = bmp.next(x, y)
				tIdx++
			}
			//countPix += T
			//tIdx += int(T)
		}
	}

	return img
}

func (bmp *BitMap) ToImageWithPalette(plt []color.Color) *ebiten.Image {
	return bmp.toImage(plt)
}

func (bmp *BitMap) next(x, y int) (int, int) {
	x++
	if x == int(bmp.GetWidth()) {
		y++
		x = 0
	}
	return x, y
}

func PixToRGBA(pix byte, plt []color.Color) color.Color {
	//return palette.Plan9[pix]
	return plt[pix]
}

/*
type _MapChunk struct {
	YJ1_FILEHEADER								// 15 bytes
	[HuffmanTreeLength*2]byte					// HuffmanTree, value of non-leaf indicates left child, value of leaf indicates the uncompressed code
	[(upperBound(HuffmanTreeLength*2) >> 4)*2]	// bitmap, indicates whether a tree node is leaf
	[BlockCount] struct {
		UncompressedLength        uint16 		// maximum 0x4000
		CompressedLength          uint16 		// including the header
		CompressedLength == 0 {
			unCompressedBlockData []byte
		} else {
			LZSSRepeatTable           [4]uint16
			LZSSOffsetCodeLengthTable [4]uint8
			LZSSRepeatCodeLengthTable [3]uint8
			CodeCountCodeLengthTable  [3]uint8
			CodeCountTable            [2]uint8
			[] {
				// one new segment with loopCount element(byte in uncompressed)
				loopCount bitStream {
					0:
						00: header->CodeCountTable[1]
						v 01-11 : next CodeCountCodeLengthTable[v - 1] bit
					1: CodeCountTable[0]
				}
				[loopCount] {
					root to leaf path in bit
				}

				// loopCount segments that each one repeats one of prev segment
				loopCount bitStream // ditto
				[loopCount] {
					count bitStream {
						00: header->LZSSRepeatTable[0]
						v 01 - 11:
							0: header->LZSSRepeatTable[v]
							1: next LZSSRepeatCodeLengthTable[temp - 1] bit
					}
					pos bitStream {
						v 00-11: next LZSSOffsetCodeLengthTable[v] bit
					}
					// repeat count times *dest = *dest - pos
				}
			}
		}
	}
}

*/

type MapChunk struct {
	data []byte
}

func (mc *MapChunk) Decompress() (*Map, error) {
	hdr := getYJ1_FILEHEADER(mc.data)
	if hdr.Signature != 0x315f4a59 { // "YJ1_"
		return nil, fmt.Errorf("")
	}
	tree_len := hdr.HuffmanTreeLength * 2
	flag := mc.data[16+tree_len:]

	root := makeHFMTree(mc.data[16:16+tree_len], flag, int(tree_len))

	var offset int
	if tree_len&0xf == 0 {
		offset = int(16 + tree_len + 2*(tree_len>>4))
	} else {
		offset = int(16 + tree_len + 2*((tree_len>>4)+1))
	}
	dst := make([]byte, hdr.UncompressedLength)
	dstPtr := 0
	for i := 0; i < int(hdr.BlockCount); i++ {
		header := get_YJ_1_BLOCKHEADER(mc.data[offset:])
		if header.CompressedLength == 0 {
			// block is not compressed, copy it directly
			offset += 4
			copy(dst[dstPtr:], mc.data[offset:offset+int(header.UncompressedLength)])
			offset += int(header.UncompressedLength)
		} else {
			offset += 24
			br := NewBitReader(mc.data[offset:])

			// read a new block
			loop := yj1_get_loop(&br, header)
			if loop == 0 {
				break
			}
			for ; loop > 0; loop-- {
				node := root
				for !node.leaf {
					if br.Read(1) == 0 {
						node = node.left
					} else {
						node = node.right
					}
				}
				dst[dstPtr] = node.value
				dstPtr++
			}

			// read a repeated block
			loop = yj1_get_loop(&br, header)
			if loop == 0 {
				break
			}
			for ; loop != 0; loop-- {
				count := yj1_get_count(&br, header)
				pos := br.Read(2)
				pos = br.Read(int(header.LZSSOffsetCodeLengthTable[pos]))
				for ; count != 0; count-- {
					dst[dstPtr] = dst[dstPtr-int(pos)]
					dstPtr++
				}
			}

			offset += int(header.CompressedLength)
		}
	}
	return &Map{dst}, nil
}

func yj1_get_count(br *BitReader, header _YJ_1_BLOCKHEADER) uint16 {
	tmp := br.Read(2)
	if tmp == 0 {
		return header.LZSSRepeatTable[0]
	} else {
		if br.Read(1) == 0 {
			return header.LZSSRepeatTable[int(tmp)]
		} else {
			return br.Read(int(header.LZSSRepeatCodeLengthTable[tmp-1]))
		}
	}
}

func yj1_get_loop(br *BitReader, header _YJ_1_BLOCKHEADER) uint16 {
	if br.Read(1) != 0 {
		return uint16(header.CodeCountTable[0])
	} else {
		tmp := br.Read(2)
		if tmp == 0 {
			return uint16(header.CodeCountTable[1])
		} else {
			return br.Read(int(header.CodeCountCodeLengthTable[tmp-1]))
		}
	}
}

func getYJ1_FILEHEADER(data []byte) _YJ1_FILEHEADER {
	// maybe use unsafe pointer is ok, but for golang i don't like do it that way
	ret := _YJ1_FILEHEADER{}
	ret.Signature = INT(binary.LittleEndian.Uint32(data[0:4]))
	ret.UncompressedLength = INT(binary.LittleEndian.Uint32(data[4:8]))
	ret.CompressedLength = INT(binary.LittleEndian.Uint32(data[8:12]))
	ret.BlockCount = uint16(binary.LittleEndian.Uint16(data[12:14]))
	ret.Unknown = uint8(data[14])
	ret.HuffmanTreeLength = uint8(data[15])
	return ret
}

func get_YJ_1_BLOCKHEADER(data []byte) _YJ_1_BLOCKHEADER {
	// ditto, maybe unsafe pointer is ok
	ret := _YJ_1_BLOCKHEADER{}
	ret.UncompressedLength = binary.LittleEndian.Uint16(data[0:2])
	ret.CompressedLength = binary.LittleEndian.Uint16(data[2:4])
	dataIdx := 4
	if ret.CompressedLength != 0 {
		for i := 0; i < len(ret.LZSSRepeatTable); i++ {
			ret.LZSSRepeatTable[i] = binary.LittleEndian.Uint16(data[dataIdx : dataIdx+2])
			dataIdx += 2
		}
		for i := range ret.LZSSOffsetCodeLengthTable {
			ret.LZSSOffsetCodeLengthTable[i] = data[dataIdx]
			dataIdx++
		}
		for i := range ret.LZSSRepeatCodeLengthTable {
			ret.LZSSRepeatCodeLengthTable[i] = data[dataIdx]
			dataIdx++
		}
		for i := range ret.CodeCountCodeLengthTable {
			ret.CodeCountCodeLengthTable[i] = data[dataIdx]
			dataIdx++
		}
		for i := range ret.CodeCountTable {
			ret.CodeCountTable[i] = data[dataIdx]
			dataIdx++
		}
	}
	return ret
}

type Map struct {
	data []byte
}

type PaletteChunk struct {
	data []byte
}

func (pc *PaletteChunk) GetPalette(night bool) ([]color.Color, error) {
	if len(pc.data) < 256*3 {
		return nil, fmt.Errorf("")
	} else if len(pc.data) < 256*3*2 {
		night = false
	}
	buf := pc.data
	ret := []color.Color{}
	offset := 0
	if night {
		offset = 3 * 256
	}
	for i := 0; i < 256; i++ {
		// if night, + 256 * 3
		r := buf[offset+i*3] << 2
		g := buf[offset+i*3+1] << 2
		b := buf[offset+i*3+2] << 2
		a := uint8(color.Opaque.A)
		ret = append(ret, color.RGBA{r, g, b, a})
	}
	return ret, nil
}

type _YJ1_FILEHEADER struct {
	Signature          INT
	UncompressedLength INT
	CompressedLength   INT
	BlockCount         uint16
	Unknown            uint8
	HuffmanTreeLength  uint8
}

type _YJ_1_BLOCKHEADER struct {
	UncompressedLength        uint16 // maximum 0x4000
	CompressedLength          uint16 // including the header
	LZSSRepeatTable           [4]uint16
	LZSSOffsetCodeLengthTable [4]uint8
	LZSSRepeatCodeLengthTable [3]uint8
	CodeCountCodeLengthTable  [3]uint8
	CodeCountTable            [2]uint8
}
