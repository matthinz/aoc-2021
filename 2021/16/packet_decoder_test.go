package d16

import (
	"strings"
	"testing"
)

func TestBitReader(t *testing.T) {
	input := "D2FE28"
	data := parseInput(strings.NewReader(input))
	r := newBitReader(data)

	version, err := r.read8(3)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if version != 6 {
		t.Fatalf("Version wrong. Expected 6, got %d", version)
	}

	typeId, err := r.read8(3)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if typeId != 4 {
		t.Fatalf("Type ID wrong. Expected 4, got %d", typeId)
	}

	next, err := r.read8(5)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if next != 0b10111 {
		t.Fatalf("Next byte wrong. expected %d, got %d", 0b10111, next)
	}
}

func TestBitReaderRead15BitValue(t *testing.T) {
	data := []byte{
		0,
		0,
		0b0101,
		0b0011,
		0b0111,
		0b0011,
	}
	r := newBitReader(data)

	expected := []uint32{
		0,
		0,
		0,
		5340,
	}

	sizes := []int{
		3,
		3,
		1,
		15,
	}

	for i := 0; i < len(expected); i++ {
		value, err := r.read32(sizes[i])
		if err != nil {
			t.Fatalf(err.Error())
		}
		if value != expected[i] {
			t.Fatalf("Read #%d failed (%d bits). Expected %d, got %d", i, sizes[i], expected[i], value)
		}

	}

}

func TestParseLiteralPacket(t *testing.T) {

	input := "D2FE28"
	data := parseInput(strings.NewReader(input))
	p := parseRootPacket(data)

	if p.Version != 6 {
		t.Fatalf("Wrong version. Expected 6, got %d", p.Version)
	}

	if p.TypeId != 4 {
		t.Fatalf("Wrong typeId. Expected 4, got %d", p.TypeId)
	}

	expectedNum := uint64(2021)
	if p.LiteralValue != expectedNum {
		t.Fatalf("Wrong literal value. Expected %d, got %d", expectedNum, p.LiteralValue)
	}

}

func TestParseOperatorPacketLengthTypeId0(t *testing.T) {

	input := "38006F45291200"
	data := parseInput(strings.NewReader(input))
	p := parseRootPacket(data)

	if p.Version != 1 {
		t.Fatalf("Wrong version. Expected 6, got %d", p.Version)
	}

	if p.TypeId != 6 {
		t.Fatalf("Wrong typeId. Expected 6, got %d", p.TypeId)
	}

	if len(p.Subpackets) != 2 {
		t.Fatalf("Wrong # of subpackets. Expected 2, got %d", len(p.Subpackets))
	}

	expectedVersions := []uint8{
		0b110,
		0b010,
	}
	expectedTypeIds := []uint8{
		0b100,
		0b100,
	}
	expectedLiteralValues := []uint64{
		0b1010,
		0b00010100,
	}

	for i := 0; i < 2; i++ {
		if p.Subpackets[i].Version != expectedVersions[i] {
			t.Fatalf("Subpacket %d has wrong version. Expected %d, got %d", i, expectedVersions[i], p.Subpackets[i].Version)
		}

		if p.Subpackets[i].TypeId != expectedTypeIds[i] {
			t.Fatalf("Subpacket %d has wrong type id. Expected %d, got %d", i, expectedTypeIds[i], p.Subpackets[i].TypeId)
		}

		if p.Subpackets[i].LiteralValue != expectedLiteralValues[i] {
			t.Fatalf("Subpacket %d has wrong literal value. Expected %d, got %d", i, expectedLiteralValues[i], p.Subpackets[i].LiteralValue)
		}

		if len(p.Subpackets[i].Subpackets) != 0 {
			t.Fatalf("Subpacket %d should not have subpackets. Expected %d, got %d", i, 0, len(p.Subpackets[i].Subpackets))
		}
	}
}

func TestParseOperatorPacketLengthTypeId1(t *testing.T) {

	input := "EE00D40C823060"
	data := parseInput(strings.NewReader(input))
	p := parseRootPacket(data)

	if p.Version != 7 {
		t.Fatalf("Wrong version. Expected 7, got %d", p.Version)
	}

	if p.TypeId != 3 {
		t.Fatalf("Wrong typeId. Expected 3, got %d", p.TypeId)
	}

	if len(p.Subpackets) != 3 {
		t.Fatalf("Wrong # of subpackets. Expected 3, got %d", len(p.Subpackets))
	}

	expectedVersions := []uint8{
		0b010,
		0b100,
		0b001,
	}
	expectedTypeIds := []uint8{
		0b100,
		0b100,
		0b100,
	}
	expectedLiteralValues := []uint64{
		0b0001,
		0b0010,
		0b0011,
	}

	for i := 0; i < 3; i++ {
		if p.Subpackets[i].Version != expectedVersions[i] {
			t.Fatalf("Subpacket %d has wrong version. Expected %d, got %d", i, expectedVersions[i], p.Subpackets[i].Version)
		}

		if p.Subpackets[i].TypeId != expectedTypeIds[i] {
			t.Fatalf("Subpacket %d has wrong type id. Expected %d, got %d", i, expectedTypeIds[i], p.Subpackets[i].TypeId)
		}

		if p.Subpackets[i].LiteralValue != expectedLiteralValues[i] {
			t.Fatalf("Subpacket %d has wrong literal value. Expected %d, got %d", i, expectedLiteralValues[i], p.Subpackets[i].LiteralValue)
		}

		if len(p.Subpackets[i].Subpackets) != 0 {
			t.Fatalf("Subpacket %d should not have subpackets. Expected %d, got %d", i, 0, len(p.Subpackets[i].Subpackets))
		}
	}
}

func TestSumVersions(t *testing.T) {

	tests := map[string]int{
		"8A004A801A8002F478":             16,
		"620080001611562C8802118E34":     12,
		"C0015000016115A2E0802F182340":   23,
		"A0016C880162017C3686B18A3D4780": 31,
	}

	for input, expected := range tests {
		data := parseInput(strings.NewReader(input))
		p := parseRootPacket(data)
		actual := p.sumVersions()
		t.Logf("expected: %d, actual: %d", expected, actual)
		if actual != expected {
			t.Fatalf("Sum of version numbers is wrong. Expected %d, got %d", expected, actual)
		}
	}

}

func TestEvaluate(t *testing.T) {
	tests := map[string]int{
		"C200B40A82":                 3,
		"04005AC33890":               54,
		"880086C3E88112":             7,
		"CE00C43D881120":             9,
		"D8005AC2A8F0":               1,
		"F600BC2D8F":                 0,
		"9C005AC2F8F0":               0,
		"9C0141080250320F1802104A08": 1,
	}

	for input, expected := range tests {

		data := parseInput(strings.NewReader(input))
		packet := parseRootPacket(data)

		actual := packet.evaluate()
		if actual != uint64(expected) {
			t.Fatalf("Expected %d, got %d", expected, actual)
		}
	}
}
