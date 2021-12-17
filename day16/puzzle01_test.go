package main

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

func TestParseLiteralPacket(t *testing.T) {

	input := "D2FE28"
	data := parseInput(strings.NewReader(input))
	p := parseRootPacket(data)

	if p.version != 6 {
		t.Fatalf("Wrong version. Expected 6, got %d", p.version)
	}

	if p.typeId != 4 {
		t.Fatalf("Wrong typeId. Expected 4, got %d", p.typeId)
	}

	expectedNum := uint32(0b011111100101)
	if p.literalValue != expectedNum {
		t.Fatalf("Wrong literal value. Expected %d, got %d", expectedNum, p.literalValue)
	}

}

func TestParseOperatorPacketLengthTypeId0(t *testing.T) {

	input := "38006F45291200"
	data := parseInput(strings.NewReader(input))
	p := parseRootPacket(data)

	if p.version != 1 {
		t.Fatalf("Wrong version. Expected 6, got %d", p.version)
	}

	if p.typeId != 6 {
		t.Fatalf("Wrong typeId. Expected 6, got %d", p.typeId)
	}

	if len(p.subpackets) != 2 {
		t.Fatalf("Wrong # of subpackets. Expected 2, got %d", len(p.subpackets))
	}

	if p.subpackets[0].literalValue != 10 {
		t.Fatalf("Subpacket 1 has wrong value. Expected %d, got %d", 10, p.subpackets[0].literalValue)
	}

	if p.subpackets[1].literalValue != 20 {
		t.Fatalf("Subpacket 3 has wrong value. Expected %d, got %d", 10, p.subpackets[1].literalValue)
	}

}

func TestSumVersion(t *testing.T) {

	tests := map[string]int{
		"8A004A801A8002F478":             16,
		"620080001611562C8802118E34":     12,
		"C0015000016115A2E0802F182340":   23,
		"A0016C880162017C3686B18A3D4780": 31,
	}

	for input, expected := range tests {
		data := parseInput(strings.NewReader(input))
		p := parseRootPacket(data)
		actual := addVersionNumbers(&p)
		if actual != expected {
			t.Fatalf("Sum of version numbers is wrong. Expected %d, got %d", expected, actual)
		}
	}

}
