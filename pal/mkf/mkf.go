package mkf

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"unsafe"
)

type (
	INT    uint32
	WORD   = uint16
	SHORT  = int16
	USHORT = uint16
	DWORD  = uint32
)

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

func (mkf *Mkf) getCompressedChunk(chunkNum INT) (CompressedChunk, error) {
	buf, err := mkf.ReadChunk(chunkNum)
	if err != nil {
		return CompressedChunk{}, err
	}
	return NewCompressedChunk(buf), nil
}

func (mkf *Mkf) GetBitMapChunk(chunkNum INT) (BitMapChunk, error) {
	buf, err := mkf.ReadChunk(chunkNum)
	if err != nil {
		return BitMapChunk{}, err
	}
	return NewBitMapChunk(buf), nil
}

/*
func (mkf *Mkf) getMgoChunk(chunkNum INT) (MgoChunk, error) {
	chunk, err := mkf.getCompressedChunk(chunkNum)
	if err != nil {
		return MgoChunk{}, err
	}
	buf, err := chunk.Decompress()
	if err != nil {
		return MgoChunk{}, err
	}
	return MgoChunk{NewBitMapChunk(buf)}, nil
}
*/

/*
func (mkf *Mkf) LoadData(chunkNum INT) ([]byte, error) {
	buf, err := mkf.ReadChunk(chunkNum)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(buf)-1; i += 2 {
		//tmp := binary.LittleEndian.Uint16(buf[i : i+2])
		p := (*uint16)(unsafe.Pointer(&buf[i]))
		*p = binary.LittleEndian.Uint16(buf[i : i+2])
	}
	return buf, nil
}
*/

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
	data      []byte
	frameSize int
}

func NewFrameChunk(data []byte) FrameChunk {
	return FrameChunk{data: data, frameSize: 320 * 200}
}

func (c *FrameChunk) SetFrameSize(frameSize int) {
	c.frameSize = frameSize
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
		return nil, fmt.Errorf("In GetFrame(), invailid frame number: %d", frameNum)
	}
	data := c.data
	offset := c.getOffset(frameNum)
	nextOffset := c.getOffset(frameNum + 1)
	fmt.Printf("%d, %d\n", offset, nextOffset)
	//if offset == 0x18444 {
	//	offset = WORD(offset)
	//}
	if nextOffset == 0 {
		return data[offset:], nil
	}
	return data[offset:nextOffset], nil
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

type CompressedChunk struct {
	data []byte
}

func NewCompressedChunk(data []byte) CompressedChunk {
	return CompressedChunk{data: data}
}

func (cc *CompressedChunk) Decompress() ([]byte, error) {
	hdr := getYJ1_FILEHEADER(cc.data)
	if hdr.Signature != 0x315f4a59 { // "YJ1_"
		return nil, fmt.Errorf("")
	}
	tree_len := int(hdr.HuffmanTreeLength) * 2
	flag := cc.data[16+tree_len:]

	root := makeHFMTree(cc.data[16:16+tree_len], flag, int(tree_len))
	root.Print()

	var offset int
	if tree_len&0xf == 0 {
		offset = int(16 + tree_len + 2*(tree_len>>4))
	} else {
		offset = int(16 + tree_len + 2*((tree_len>>4)+1))
	}
	dst := make([]byte, hdr.UncompressedLength)
	dstPtr := 0
	for i := 0; i < int(hdr.BlockCount); i++ {
		header := get_YJ_1_BLOCKHEADER(cc.data[offset:])
		if header.CompressedLength == 0 {
			// block is not compressed, copy it directly
			offset += 4 // block header len
			copy(dst[dstPtr:], cc.data[offset:offset+int(header.UncompressedLength)])
			offset += int(header.UncompressedLength)
		} else {
			offset += 24 // block header len
			br := NewBitReader(cc.data[offset:])

			for j := 0; ; j++ {
				//fmt.Println(i, j)
				if i == 1 && (j == 60) {
					fmt.Println()
				}
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
			}
			offset += int(header.CompressedLength) - 24
		}
	}
	return dst, nil
}

func yj1_get_count(br *BitReader, header _YJ_1_BLOCKHEADER) uint16 {
	tmp := br.Read(2)
	var ret uint16
	if tmp == 0 {
		ret = header.LZSSRepeatTable[0]
	} else {
		if br.Read(1) == 0 {
			ret = header.LZSSRepeatTable[int(tmp)]
		} else {
			ret = br.Read(int(header.LZSSRepeatCodeLengthTable[tmp-1]))
		}
	}
	//if ret > 100 {
	//	ret += 5000
	//}
	return ret
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

type PlaneChunk struct {
	data    []byte
	eleSize int
}

func NewPlaneChunk(data []byte, eleSize int) PlaneChunk {
	for i := 0; i < len(data)-1; i += 2 {
		//tmp := binary.LittleEndian.Uint16(buf[i : i+2])
		p := (*uint16)(unsafe.Pointer(&data[i]))
		*p = binary.LittleEndian.Uint16(data[i : i+2])
	}
	return PlaneChunk{data: data, eleSize: eleSize}
}

func (pc *PlaneChunk) Len() int {
	return len(pc.data) / pc.eleSize
}

func (pc *PlaneChunk) Get(idx int) unsafe.Pointer {
	i := idx * int(pc.eleSize)
	return unsafe.Pointer(&pc.data[i])
}
