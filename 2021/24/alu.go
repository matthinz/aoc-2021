package d24

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/matthinz/aoc-golang"
)

type aluOperationKind int

const (
	OpInput aluOperationKind = iota
	OpAdd
	OpMultiply
	OpDivide
	OpModulo
	OpEquality
)

type aluRegister int

const (
	RegisterW aluRegister = iota
	RegisterX
	RegisterY
	RegisterZ
)

type alu struct {
	w, x, y, z int
	inputs     []int
}

type literalArg struct {
	value int
}

type registerArg struct {
	register aluRegister
}

type aluOperation struct {
	kind aluOperationKind
	lhs  registerArg
	rhs  aluOperationArg
}

type aluOperationArg interface {
	Value(registers *alu) int
}

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(24, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {
	return ""

}

func Puzzle2(r io.Reader, l *log.Logger) string {
	return ""
}

////////////////////////////////////////////////////////////////////////////////
// brute force

func bruteForcePuzzle1(r io.Reader, l *log.Logger) string {

	const NumberLength = 14

	ops := parseInput(r)

	digits := make([]int, NumberLength)
	for i := 0; i < NumberLength; i++ {
		digits[i] = 9
	}

	var processedCount uint

	for {

		processedCount++
		if processedCount%1000000 == 0 {
			l.Printf("Processed %d numbers", processedCount)
		}

		a := executeAll(digits, ops)

		if a.z == 0 {
			result := ""
			for _, digit := range digits {
				result = result + strconv.FormatInt(int64(digit), 10)
			}
			return result
		}

		if digits[NumberLength-1] > 1 {
			digits[NumberLength-1]--
			continue
		}

		// we need to take from successive positions
		didTake := false
		for i := NumberLength - 2; i >= 0; i-- {
			if digits[i] > 1 {
				digits[i]--
				for j := i + 1; j < NumberLength; j++ {
					digits[j] = 9
				}
				didTake = true
				break
			}
		}

		if !didTake {
			break
		}
	}

	panic("Could not find valid number")

}

////////////////////////////////////////////////////////////////////////////////
// execute

func executeAll(inputs []int, ops []aluOperation) *alu {
	result := &alu{
		inputs: inputs,
	}
	for i := range ops {
		result = execute(result, &ops[i])
	}
	return result
}

func execute(a *alu, op *aluOperation) *alu {

	next := alu{
		inputs: a.inputs,
		w:      a.w,
		x:      a.x,
		y:      a.y,
		z:      a.z,
	}

	var result int

	switch op.kind {
	case OpAdd:
		result = op.lhs.Value(a) + op.rhs.Value(a)

	case OpDivide:
		result = op.lhs.Value(a) / op.rhs.Value(a)

	case OpEquality:
		if op.lhs.Value(a) == op.rhs.Value(a) {
			result = 1
		} else {
			result = 0
		}

	case OpInput:
		result = a.inputs[0]
		next.inputs = next.inputs[1:]

	case OpModulo:
		result = op.lhs.Value(a) % op.rhs.Value(a)

	case OpMultiply:
		result = op.lhs.Value(a) * op.rhs.Value(a)

	default:
		panic(fmt.Sprintf("Invalid op code: %v", op.kind))
	}

	next.setRegister(op.lhs.register, result)

	return &next
}

////////////////////////////////////////////////////////////////////////////////
// various methods

func (a *alu) setRegister(register aluRegister, value int) {
	switch register {
	case RegisterW:
		a.w = value
	case RegisterX:
		a.x = value
	case RegisterY:
		a.y = value
	case RegisterZ:
		a.z = value
	default:
		panic(fmt.Sprintf("Invalid register: %v", register))
	}
}

func (a literalArg) Value(registers *alu) int {
	return a.value
}

func (a registerArg) Value(registers *alu) int {
	switch a.register {
	case RegisterW:
		return registers.w
	case RegisterX:
		return registers.x
	case RegisterY:
		return registers.y
	case RegisterZ:
		return registers.z
	default:
		panic(fmt.Sprintf("Invalid register %v", a.register))
	}
}

////////////////////////////////////////////////////////////////////////////////
// parseInput

func parseInput(r io.Reader) []aluOperation {

	kinds := map[string]aluOperationKind{
		"inp": OpInput,
		"add": OpAdd,
		"mul": OpMultiply,
		"div": OpDivide,
		"mod": OpModulo,
		"eql": OpEquality,
	}

	registers := map[string]aluRegister{
		"w": RegisterW,
		"x": RegisterX,
		"y": RegisterY,
		"z": RegisterZ,
	}

	var ops []aluOperation

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) != 2 && len(parts) != 3 {
			continue
		}

		kind, ok := kinds[parts[0]]
		if !ok {
			panic(fmt.Sprintf("Invalid kind: %s", parts[0]))
		}

		register, ok := registers[parts[1]]
		if !ok {
			panic(fmt.Sprintf("Invalid register: %s", parts[1]))
		}

		var rhs aluOperationArg

		if len(parts) > 2 {

			value, err := strconv.ParseInt(parts[2], 10, 32)
			if err == nil {
				// value is a literal
				rhs = literalArg{
					value: int(value),
				}
			} else {
				// rhs is a register
				register, ok := registers[parts[2]]
				if !ok {
					panic(fmt.Sprintf("Invalid register: %s", parts[2]))
				}
				rhs = registerArg{
					register: register,
				}
			}
		}

		op := aluOperation{
			kind: kind,
			lhs: registerArg{
				register: register,
			},
			rhs: rhs,
		}

		ops = append(ops, op)
	}

	return ops
}
