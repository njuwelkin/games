package main

type BitReader struct {
	data   []byte
	bitIdx int
}

func NewBitReader(data []byte) BitReader {
	return BitReader{data: data}
}

func (br *BitReader) Read(count int) uint16 {
	tmp := (br.bitIdx >> 4) << 1 // read two neibor bytes one time, tmp indicates index of first one in data array
	bptr := br.bitIdx & 0xf      // offset in two bytes
	br.bitIdx += count
	if count > 16-bptr {
		count = count + bptr - 16
		mask := uint16(0xffff) >> bptr
		return (((uint16(br.data[tmp]) | (uint16(br.data[tmp+1]) << 8)) & mask) << count) |
			((uint16(br.data[tmp+2]) | (uint16(br.data[tmp]+3) << 8)) >> (16 - count))
	} else {
		/*
			a := uint16(br.data[tmp])
			b := uint16(br.data[tmp+1]) << 8
			c := a | b
			d := c << uint16(bptr)
			fmt.Printf("d: %b\n", d)
			e := d >> (16 - count)
			fmt.Printf("e: %b\n", e)
			return e
		*/
		return ((uint16(br.data[tmp]) | (uint16(br.data[tmp+1]) << 8)) << bptr) >> (16 - count)
	}
}

func (br *BitReader) Reset() {
	br.bitIdx = 0
}
