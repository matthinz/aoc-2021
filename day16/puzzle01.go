package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// bitReader lets you read through a slice of bytes an arbitrary number of
// bits at a time (< 8)
type bitReader struct {
	data      []byte
	pos       int
	bitOffset int
}

type packet struct {
	version      uint8
	typeId       uint8
	literalValue uint32
	subpackets   []packet
}

const LiteralPacketTypeId = 4

func main() {

	data := parseInput(os.Stdin)

	packet := parseRootPacket(data)

	fmt.Println(addVersionNumbers(&packet))

}

func addVersionNumbers(p *packet) int {
	result := int(p.version)

	for i := range p.subpackets {
		subpacket := &p.subpackets[i]
		result += addVersionNumbers(subpacket)
	}

	return result
}

////////////////////////////////////////////////////////////////////////////////
// parseRootPacket

func parseRootPacket(data []byte) packet {

	r := newBitReader(data)

	p, _ := parsePacket(r)

	return *p
}

// parses a single packet out of r and returns it along with its bit length
func parsePacket(r *bitReader) (*packet, int) {

	version, err := r.read8(3)
	if err == io.EOF {
		return nil, 0
	}

	typeId, err := r.read8(3)
	if err == io.EOF {
		return nil, 0
	}

	headerLength := 3 + 3

	if typeId == LiteralPacketTypeId {
		literalPacket, length := parseLiteralPacket(version, typeId, r)
		return (&literalPacket), length + headerLength
	} else {
		operatorPacket, length := parseOperatorPacket(version, typeId, r)
		return (&operatorPacket), length + headerLength
	}

}

func parseLiteralPacket(version uint8, typeId uint8, r *bitReader) (packet, int) {
	const wordSize = 4

	p := packet{
		version: version,
		typeId:  typeId,
	}

	length := 0

	for i := 0; i < 64/wordSize; i++ {
		// make room
		p.literalValue = p.literalValue << wordSize

		chunk, err := r.read8(5)
		if err != nil {
			panic(err)
		}
		length += 5

		mask := uint8(0b11101111)
		p.literalValue = p.literalValue | uint32(chunk&mask)

		isLastChunk := chunk&mask == chunk

		if isLastChunk {
			break
		}
	}

	return p, length
}

func parseOperatorPacket(version uint8, typeId uint8, r *bitReader) (packet, int) {

	length := 0

	lengthTypeId, err := r.read8(1)
	if err != nil {
		panic(err)
	}

	length += 1

	p := packet{
		version: version,
		typeId:  typeId,
	}

	if lengthTypeId == 0 {
		// If the length type ID is 0, then the next 15 bits are a number that
		// represents the total length in bits of the sub-packets contained by
		// this packet.
		expectedLengthInBits, err := r.read32(15)
		if err != nil {
			panic(err)
		}
		length += 15

		subpacketBitsProcessed := 0
		for subpacketBitsProcessed < int(expectedLengthInBits) {
			subpacket, subpacketLength := parsePacket(r)
			length += subpacketLength
			subpacketBitsProcessed += subpacketLength
			p.subpackets = append(p.subpackets, *subpacket)
		}
	} else if lengthTypeId == 1 {

		// If the length type ID is 1, then the next 11 bits are a number that
		// represents the number of sub-packets immediately contained by this packet.
		numberOfSubpackets, err := r.read32(11)
		if err != nil {
			panic(err)
		}
		length += 11

		for i := 0; i < int(numberOfSubpackets); i++ {
			subpacket, subpacketLength := parsePacket(r)
			if subpacket == nil {
				panic(fmt.Sprintf("Got nil subpacket when parsing %d / %d", i+1, numberOfSubpackets))
			}
			length += subpacketLength
			p.subpackets = append(p.subpackets, *subpacket)
		}

	} else {
		panic("invalid length type id")
	}

	return p, length

}

////////////////////////////////////////////////////////////////////////////////
// bitReader implementation

func newBitReader(data []byte) *bitReader {
	return &bitReader{data, 0, 0}
}

func (b *bitReader) atEnd() bool {
	return b.pos >= len(b.data)
}

func (b *bitReader) read8(bits int) (uint8, error) {
	if bits > 8 || bits < 0 {
		return 0, fmt.Errorf("read8 received bad number of bits: %d", bits)
	}
	value, err := b.read32(bits)
	return uint8(value), err
}

// reads up to 32 bits off of <br>
func (b *bitReader) read32(bits int) (uint32, error) {

	if bits > 32 || bits < 0 {
		return 0, fmt.Errorf("read32 received bad number of bits: %d", bits)
	}

	// NOTES
	// - Even though we have 8 bits in each byte of data, we only use 4 of them

	const wordSize = 4

	var result uint32
	bitsRead := 0

	for bitsRead < bits {

		currentByte := byte(0)
		if b.pos < len(b.data) {
			currentByte = b.data[b.pos]
		}

		bitsToRead := bits - bitsRead
		bitsAvailable := wordSize - b.bitOffset
		if bitsAvailable < bitsToRead {
			bitsToRead = bitsAvailable
		}

		// mask off the previously read bits
		mask := byte(0xFF) >> (8 - wordSize) >> b.bitOffset
		currentByte = currentByte & mask

		// Shift to remove the bits we don't care about
		bitsLeft := wordSize - (b.bitOffset + bitsToRead)
		currentByte = currentByte >> byte(bitsLeft)

		// Shift back to make room for any more bits we'll need to set
		bitsInFutureBytes := bits - (bitsRead + bitsToRead)
		currentByte = currentByte << byte(bitsInFutureBytes)

		result = result | uint32(currentByte)

		bitsRead += bitsToRead
		b.bitOffset += bitsToRead

		if b.bitOffset >= wordSize {
			b.pos++
			b.bitOffset = 0
		}
	}

	if b.pos >= len(b.data) {
		return result, io.EOF
	}

	return result, nil
}

////////////////////////////////////////////////////////////////////////////////

func parseInput(r io.Reader) []uint8 {
	result := make([]uint8, 0)

	buffer := make([]byte, 16)
	for {
		ct, err := r.Read(buffer)

		if err != nil && err != io.EOF {
			panic(err)
		}
		if ct == 0 {
			break
		}
		for i := 0; i < ct; i++ {
			c := rune(buffer[i])
			if c == '\n' {
				continue
			}
			value, err := strconv.ParseInt(string(c), 16, 8)
			if err != nil {
				panic(err)
			}
			result = append(result, uint8(value))
		}
	}

	return result
}

////////////////////////////////////////////////////////////////////////////////

func format8Bits(value uint8) string {
	b := strings.Builder{}

	for i := 7; i >= 0; i-- {
		mask := byte(1) << i
		if value&mask == mask {
			b.WriteString("1")
		} else {
			b.WriteString("0")
		}
	}
	return b.String()
}

func format32Bits(value uint32) string {
	b := strings.Builder{}

	for i := 32; i >= 0; i-- {
		mask := uint32(1) << i
		if value&mask == mask {
			b.WriteString("1")
		} else {
			b.WriteString("0")
		}
	}
	return b.String()
}
