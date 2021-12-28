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
	debug     bool
}

type Packet struct {
	Version      uint8
	TypeId       uint8
	LiteralValue uint64
	Subpackets   []Packet
}

const LiteralPacketTypeId = 4

func main() {

	data := parseInput(os.Stdin)

	packet := parseRootPacket(data)

	fmt.Println(packet.sumVersions())

	fmt.Println(packet.evaluate())

}

func printPacket(p *Packet, prefix string) {

	fmt.Printf("%sv: %s (%d)\n", prefix, format4Bits(p.Version), p.Version)
	fmt.Printf("%st: %s (%d)\n", prefix, format4Bits(p.TypeId), p.TypeId)
	fmt.Printf("%svsum: %d\n", prefix, p.sumVersions())

	if p.TypeId == 4 {
		fmt.Printf("%svalue: %d\n", prefix, p.LiteralValue)
	} else {
		fmt.Printf("%ssp:\n", prefix)
		for i := range p.Subpackets {
			printPacket(&p.Subpackets[i], prefix+"  ")
		}
	}

}

func (p *Packet) evaluate() uint64 {

	switch p.TypeId {
	case 0: // SUM

		var result uint64
		for _, sp := range p.Subpackets {
			result += sp.evaluate()
		}
		return result

	case 1: // PRODUCT

		result := uint64(1)
		for _, sp := range p.Subpackets {
			result *= sp.evaluate()
		}
		return result

	case 2: // MINIMUM

		var result *uint64
		for _, sp := range p.Subpackets {
			val := sp.evaluate()
			if result == nil || val < (*result) {
				result = &val
			}
		}
		return *result

	case 3: // MAXIMUM

		var result *uint64
		for _, sp := range p.Subpackets {
			val := sp.evaluate()
			if result == nil || val > (*result) {
				result = &val
			}
		}
		return *result

	case 4: // LITERAL
		return p.LiteralValue

	case 5: // GREATER THAN

		if len(p.Subpackets) != 2 {
			panic("packet type 5 should always have 2 subpackets")
		}
		if p.Subpackets[0].evaluate() > p.Subpackets[1].evaluate() {
			return 1
		} else {
			return 0
		}

	case 6: // LESS THAN

		if len(p.Subpackets) != 2 {
			panic("packet type 6 should always have 2 subpackets")
		}
		if p.Subpackets[0].evaluate() < p.Subpackets[1].evaluate() {
			return 1
		} else {
			return 0
		}

	case 7: // EQUAL TO

		if len(p.Subpackets) != 2 {
			panic("packet type 5 should always have 2 subpackets")
		}
		if p.Subpackets[0].evaluate() == p.Subpackets[1].evaluate() {
			return 1
		} else {
			return 0
		}

	default:
		panic("invalid packet type!")

	}

}

func (p *Packet) sumVersions() int {
	var result int

	result += int(p.Version)

	for i := range p.Subpackets {
		result += p.Subpackets[i].sumVersions()
	}

	return result
}

////////////////////////////////////////////////////////////////////////////////
// parseRootPacket

func parseRootPacket(data []byte) Packet {

	r := newBitReader(data)
	r.debug = true

	p, _ := parsePacket(r)

	return *p
}

// parses a single packet out of r and returns it along with its bit length
func parsePacket(r *bitReader) (*Packet, int) {

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

func parseLiteralPacket(version uint8, typeId uint8, r *bitReader) (Packet, int) {
	const wordSize = 4

	p := Packet{
		Version: version,
		TypeId:  typeId,
	}

	length := 0
	hitLastChunk := false

	for i := 0; i < 64/wordSize; i++ {
		// make room
		p.LiteralValue = p.LiteralValue << wordSize

		chunk, err := r.read8(5)
		if err != nil && err != io.EOF {
			panic(err)
		}
		length += 5

		mask := uint8(0b11101111)
		p.LiteralValue = p.LiteralValue | uint64(chunk&mask)

		hitLastChunk = chunk&mask == chunk

		if hitLastChunk {
			break
		}
	}

	if !hitLastChunk {
		panic("did not find last chunk -- might have a literal > 64 bits")
	}

	return p, length
}

func parseOperatorPacket(version uint8, typeId uint8, r *bitReader) (Packet, int) {

	length := 0

	lengthTypeId, err := r.read8(1)
	if err != nil {
		panic(err)
	}

	length += 1

	p := Packet{
		Version: version,
		TypeId:  typeId,
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
			p.Subpackets = append(p.Subpackets, *subpacket)
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
			p.Subpackets = append(p.Subpackets, *subpacket)
		}

	} else {
		panic("invalid length type id")
	}

	return p, length

}

////////////////////////////////////////////////////////////////////////////////
// bitReader implementation

func newBitReader(data []byte) *bitReader {
	return &bitReader{data, 0, 0, false}
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

		current := uint32(0)
		if b.pos < len(b.data) {
			current = uint32(b.data[b.pos])
		}

		bitsToRead := bits - bitsRead
		bitsAvailable := wordSize - b.bitOffset
		if bitsAvailable < bitsToRead {
			bitsToRead = bitsAvailable
		}

		// mask off the previously read bits
		mask := uint32(0xFF) >> (8 - wordSize) >> b.bitOffset
		current = current & mask

		// Shift to remove the bits we don't care about
		bitsLeft := wordSize - (b.bitOffset + bitsToRead)
		current = current >> byte(bitsLeft)

		// Shift back to make room for any more bits we'll need to set
		bitsInFutureBytes := bits - (bitsRead + bitsToRead)
		current = current << byte(bitsInFutureBytes)

		result = result | current

		bitsRead += bitsToRead
		b.bitOffset += bitsToRead

		if b.bitOffset >= wordSize {
			b.pos++
			b.bitOffset = 0
		}
	}

	if b.debug {
		nice := format32Bits(result)
		fmt.Fprintf(os.Stderr, "READ %s\n", nice[32-bits:32])
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

func format4Bits(value uint8) string {
	b := strings.Builder{}

	for i := 3; i >= 0; i-- {
		mask := byte(1) << i
		if value&mask == mask {
			b.WriteString("1")
		} else {
			b.WriteString("0")
		}
	}
	return b.String()
}

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

	for i := 32; i > 0; i-- {
		mask := uint32(1) << i
		if value&mask == mask {
			b.WriteString("1")
		} else {
			b.WriteString("0")
		}
	}
	return b.String()
}
