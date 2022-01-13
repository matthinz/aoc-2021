package d24

import "fmt"

type Registers struct {
	w Expression
	x Expression
	y Expression
	z Expression
}

func NewRegisters() *Registers {
	return &Registers{
		w: zeroLiteral,
		x: zeroLiteral,
		y: zeroLiteral,
		z: zeroLiteral,
	}
}

func (r *Registers) set(name string, value Expression) {
	switch name {
	case "w":
		r.w = value
	case "x":
		r.x = value
	case "y":
		r.y = value
	case "z":
		r.z = value
	default:
		panic(fmt.Sprintf("Invalid register: %s", name))
	}
}

func (r *Registers) get(name string) Expression {
	switch name {
	case "w":
		return r.w
	case "x":
		return r.x
	case "y":
		return r.y
	case "z":
		return r.z
	default:
		panic(fmt.Sprintf("Invalid register: %s", name))
	}

}
