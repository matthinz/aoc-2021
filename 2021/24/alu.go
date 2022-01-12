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

type inputExpression struct {
	index int
}

type literalExpression struct {
	value int
}

type binaryExpression struct {
	lhs          Expression
	rhs          Expression
	isSimplified bool
	outputRange  *[2]int
}

type addExpression struct {
	binaryExpression
}

type divideExpression struct {
	binaryExpression
}

type equalsExpression struct {
	binaryExpression
}

type moduloExpression struct {
	binaryExpression
	cachedRange *[2]int
}

type multiplyExpression struct {
	binaryExpression
}

type namedExpression struct {
	name string
	expr Expression
}

type registers struct {
	w Expression
	x Expression
	y Expression
	z Expression
}

type inputPreference int

const (
	PreferFirstInput inputPreference = iota
	PreferSecondInput
)

// decider takes two sets of inputs and decides which one to use.
type decider func(a, b map[int]int) (map[int]int, error)

// expression is a generic interface encapsulating an expression that can
// be evaluated with an ALU
type Expression interface {
	Evaluate(inputs []int) int
	// Returns a set of inputs that will make this expression evaluate to <target>.
	// <d> is a function that, given two potential sets of inputs, returns the one that should be preferred.
	FindInputs(target int, d decider) (map[int]int, error)
	Range() IntRange
	Simplify() Expression
	String() string
}

type IntRange struct {
	min, max int
	step     int
}

const MinInputValue = 1

const MaxInputValue = 9

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(24, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {

	const Digits = 14

	reg := parseInput(r)

	inputs, err := reg.z.FindInputs(0, func(a, b map[int]int) (map[int]int, error) {
		panic("NOT IMPLEMENTED")
	})

	if err != nil {
		panic(err)
	}

	digits := make([]int, Digits)

	for i := 0; i < Digits; i++ {
		digit, isSet := inputs[i]
		if isSet {
			digits[i] = digit
		} else {
			digits[i] = 9
		}
	}

	l.Println(digits)

	result := reg.z.Evaluate(digits)
	l.Println(result)

	return ""
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	return ""
}

////////////////////////////////////////////////////////////////////////////////
// range

func (r IntRange) EqualsInt(value int) bool {
	return r.min == r.max && r.min == value
}

func (r IntRange) Len() int {
	return (r.max - r.min) + 1
}

func (r IntRange) Translate(other IntRange, t func(a, b int) int) (IntRange, error) {
	if r.step != other.step {
		return IntRange{}, fmt.Errorf("Can't combine ranges with different steps (%d and %d)", r.step, other.step)
	}
	return IntRange{
		min:  t(r.min, other.min),
		max:  t(r.max, other.max),
		step: r.step,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////
// expressions

func (e *addExpression) Evaluate(inputs []int) int {
	return e.lhs.Evaluate(inputs) + e.rhs.Evaluate(inputs)
}

func (e *addExpression) FindInputs(target int, d decider) (map[int]int, error) {
	return findInputsForBinaryExpression(
		&e.binaryExpression,
		target,
		func(lhsValue int, rhsRange IntRange) ([]int, error) {
			rhsValue := target - lhsValue
			if rhsValue < rhsRange.min || rhsValue > rhsRange.max {
				return []int{}, nil
			}
			return []int{rhsValue}, nil
		},
		d,
	)
}

func (e *addExpression) Range() IntRange {
	lhsRange := e.lhs.Range()
	rhsRange := e.rhs.Range()

	return IntRange{
		min:  lhsRange.min + rhsRange.min,
		max:  lhsRange.max + rhsRange.max,
		step: lhsRange.step * rhsRange.step,
	}
}

func (e *addExpression) Simplify() Expression {
	if e.binaryExpression.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

	// if both ranges are zero length, we are adding two literals
	if lhsRange.Len() == 0 && rhsRange.Len() == 0 {
		return &literalExpression{
			value: lhsRange.min + rhsRange.min,
		}
	}

	// if either range is zero, use the other
	if lhsRange.EqualsInt(0) {
		return rhs
	}

	if rhsRange.EqualsInt(0) {
		return lhs
	}

	return &addExpression{
		binaryExpression: binaryExpression{
			lhs:          lhs,
			rhs:          rhs,
			isSimplified: true,
		},
	}
}

func (e *addExpression) String() string {
	rhsRange := e.rhs.Range()
	if rhsRange.min < 0 && rhsRange.max < 0 {
		return fmt.Sprintf("(%s - %s)", e.lhs.String(), strings.Replace(e.rhs.String(), "-", "", 1))
	} else {
		return fmt.Sprintf("(%s + %s)", e.lhs.String(), e.rhs.String())
	}
}

func (e *divideExpression) Evaluate(inputs []int) int {
	return e.lhs.Evaluate(inputs) / e.rhs.Evaluate(inputs)
}

func (e *divideExpression) FindInputs(target int, d decider) (map[int]int, error) {
	return findInputsForBinaryExpression(
		&e.binaryExpression,
		target,
		func(dividend int, divisorRange IntRange) ([]int, error) {

			// 4th grade math recap: dividend / divisor = target
			// here we return potential divisors between min and max that will equal target

			if target == 0 {
				// When target == 0, divisor can't affect the result, except when it
				// can. We're doing integer division, so a large enough divisor *could* get us to zero
				panic("NOT IMPLEMENTED")
			}

			// dividend = divisor * target
			// divisor = dividend / target
			divisor := dividend / target

			if divisor < divisorRange.min || divisor > divisorRange.max {
				return []int{}, nil
			}

			if dividend/divisor != target {
				return []int{}, nil
			}

			return []int{divisor}, nil
		},
		d,
	)
}

func (e *divideExpression) Range() IntRange {
	lhsRange := e.lhs.Range()
	rhsRange := e.rhs.Range()
	return IntRange{
		min: lhsRange.min / rhsRange.max,
		max: lhsRange.max / rhsRange.min,
	}
}

func (e *divideExpression) Simplify() Expression {
	if e.binaryExpression.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

	// if both ranges are zero-length, we can sub in a literal
	if lhsRange.Len() == 0 && rhsRange.Len() == 0 {
		return &literalExpression{
			value: lhsRange.min / rhsRange.min,
		}
	}

	// if left value is zero, this will eval to zero
	if lhsRange.EqualsInt(0) {
		return &literalExpression{
			value: 0,
		}
	}

	// if right value is 1, this will eval to lhs
	if rhsRange.EqualsInt(1) {
		return lhs
	}

	return &divideExpression{
		binaryExpression: binaryExpression{
			lhs:          lhs,
			rhs:          rhs,
			isSimplified: true,
		},
	}
}

func (e *divideExpression) String() string {
	return fmt.Sprintf("(%s / %s)", e.lhs.String(), e.rhs.String())
}

func (e *equalsExpression) Evaluate(inputs []int) int {
	lhsValue := e.lhs.Evaluate(inputs)
	rhsValue := e.rhs.Evaluate(inputs)
	if lhsValue == rhsValue {
		return 1
	} else {
		return 0
	}
}

func (e *equalsExpression) FindInputs(target int, d decider) (map[int]int, error) {
	if target != 0 && target != 1 {
		return nil, fmt.Errorf("equalsExpression can't seek a target other than 0 or 1 (got %d)", target)
	}

	return findInputsForBinaryExpression(
		&e.binaryExpression,
		target,
		func(lhsValue int, rhsRange IntRange) ([]int, error) {
			if target == 0 {
				// We must find *any* rhsValue that does not equal lhsValue
				if rhsRange.min == rhsRange.max {
					if rhsRange.min != lhsValue {
						return []int{min}, nil
					}
					return []int{}, nil
				}

				panic("NOT IMPLEMENTED")
			}

			// we must find an value that equals lhsValue
			if lhsValue < rhsRange.min || lhsValue > rhsRange.max {
				// not in the range of possible values
				return []int{}, nil
			}
			return []int{lhsValue}, nil
		},
		d,
	)
}

func (e *equalsExpression) Range() IntRange {
	return IntRange{0, 1, 1}
}

func (e *equalsExpression) Simplify() Expression {
	if e.binaryExpression.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

	// if all elements of both ranges are equal, we are comparing two equal values
	if lhsRange.EqualsRange(rhsRange) {
		return &literalExpression{
			value: 1,
		}
	}

	// if the ranges of each side of the comparison will never intersect,
	// then we can always return "0" for this expression

	if !lhsRange.IntersectsRange(rhsRange) {
		return &literalExpression{
			value: 0,
		}
	}

	return &equalsExpression{
		binaryExpression: binaryExpression{
			lhs:          lhs,
			rhs:          rhs,
			isSimplified: true,
		},
	}
}

func (e *equalsExpression) String() string {
	return fmt.Sprintf("(%s == %s ? 1 : 0)", e.lhs.String(), e.rhs.String())
}

func (e *inputExpression) Evaluate(inputs []int) int {
	return inputs[e.index]
}

func (e *inputExpression) FindInputs(target int, d decider) (map[int]int, error) {
	if target < 0 || target > 9 {
		return nil, fmt.Errorf("inputExpression can't seek a target not in the range 1-9 (got %d)", target)
	}

	m := make(map[int]int, 1)
	m[e.index] = target
	return m, nil
}

func (e *inputExpression) Range() IntRange {
	return IntRange{
		min:  1,
		max:  9,
		step: 1,
	}
}

func (e *inputExpression) Simplify() Expression {
	return e
}

func (e *inputExpression) String() string {
	return fmt.Sprintf("i%d", e.index)
}

func (e *literalExpression) Evaluate(inputs []int) int {
	return e.value
}

func (e *literalExpression) FindInputs(target int, d decider) (map[int]int, error) {
	if e.value != target {
		return nil, fmt.Errorf("literalValue %d can't seek target value %d", e.value, target)
	}
	// no inputs can affect this expression's value
	return map[int]int{}, nil
}

func (e *literalExpression) Range() IntRange {
	return IntRange{
		min:  e.value,
		max:  e.value,
		step: 1,
	}
}

func (e *literalExpression) Simplify() Expression {
	return e
}

func (e *literalExpression) String() string {
	return strconv.FormatInt(int64(e.value), 10)
}

func (e *moduloExpression) Evaluate(inputs []int) int {
	return e.lhs.Evaluate(inputs) % e.rhs.Evaluate(inputs)
}

func (e *moduloExpression) FindInputs(target int, d decider) (map[int]int, error) {
	return findInputsForBinaryExpression(
		&e.binaryExpression,
		target,
		func(lhsValue int, rhsRange IntRange) ([]int, error) {
			// need to find rhsValues such that lhsValue % rhsValue = target
			// these would be factors of lhsValue - target
			// TODO: smarter way

			const MaxValues = 1000
			var result []int
			for i := rhsRange.min; i <= rhsRange.max; i += rhsRange.step {
				if lhsValue%i == target {
					result = append(result, i)
					if len(result) >= MaxValues {
						return nil, fmt.Errorf("Too many values between %d - %d where %d %% x == %d", rhsRange.min, rhsRange.max, lhsValue, target)
					}
				}
			}

			return result, nil
		},
		d,
	)
}

func (e *moduloExpression) Range() IntRange {

	// modulo ranges are hard

	panic("NOT IMPLEMENTED")

}

func (e *moduloExpression) Simplify() Expression {
	if e.binaryExpression.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

	// if lhs is 0, we can resolve to zero
	if lhsRange.EqualsInt(0) {
		return &literalExpression{
			value: 0,
		}
	}

	// if both ranges are zero-length, we can evaluate to a literal
	if lhsRange.Len() == 0 && rhsRange.Len() == 0 {
		return &literalExpression{
			value: lhsRange[0] % rhsRange[0],
		}
	}

	// if lhs is zero-length and *less than* the rhs range, we can eval to a literal
	if lhsRange.Len() == 0 && lhsRange.min < rhsRange.min {
		return &literalExpression{
			value: lhsRange[0],
		}
	}

	return &moduloExpression{
		binaryExpression: binaryExpression{
			lhs:          lhs,
			rhs:          rhs,
			isSimplified: true,
		},
	}
}

func (e *moduloExpression) String() string {
	return fmt.Sprintf("(%s %% %s)", e.lhs.String(), e.rhs.String())
}

func (e *multiplyExpression) Evaluate(inputs []int) int {
	return e.lhs.Evaluate(inputs) * e.rhs.Evaluate(inputs)
}

func (e *multiplyExpression) FindInputs(target int, d decider) (map[int]int, error) {
	return findInputsForBinaryExpression(
		&e.binaryExpression,
		target,
		func(lhsValue, min, max int) ([]int, error) {
			if target == 0 {
				if lhsValue != 0 {
					// rhsValue *must* be zero
					if min <= 0 && max >= 0 {
						return []int{0}, nil
					}
				}

				// lhsValue is zero, so rhsValue can be literally *any* number
				if min == max {
					return []int{min}, nil
				}

				return nil, fmt.Errorf("Too many values between %d - %d such that %d * x = %d", min, max, lhsValue, target)
			}

			if target == lhsValue {
				if min <= 1 && max >= 1 {
					return []int{1}, nil
				} else {
					return []int{}, nil
				}
			}

			const MaxValues = 1000
			var result []int

			for i := min; i <= max; i++ {
				if lhsValue*i == target {
					result = append(result, i)
					if len(result) >= MaxValues {
						return nil, fmt.Errorf("Too many values multiply by %d to get %d", lhsValue, target)
					}
				}
			}

			return result, nil
		},
		d,
	)
}

func (e *multiplyExpression) Range() IntRange {
	return findBinaryExpressionRange(
		&e.binaryExpression,
		func(lhs, rhs int) (int, error) {
			return lhs * rhs, nil
		},
	)
}

func (e *multiplyExpression) Simplify() Expression {
	if e.binaryExpression.isSimplified {
		return e
	}

	lhs := e.lhs.Simplify()
	rhs := e.rhs.Simplify()

	lhsRange := lhs.Range()
	rhsRange := rhs.Range()

	// if both ranges are zero-length, we are doing literal multiplication
	if lhsRange.Len() == 0 && rhsRange.Len() == 0 {
		return &literalExpression{
			value: lhsRange.min * rhsRange.min,
		}
	}

	// if either range is just "0", we'll evaluate to 0
	if (lhsRange.EqualsInt(0)) || (rhsRange.EqualsInt(0)) {
		return &literalExpression{
			value: 0,
		}
	}

	// if either range is just "1", we evaluate to the other
	if lhsRange.EqualsInt(1) {
		return rhs
	}

	if rhsRange.EqualsInt(1) {
		return lhs
	}

	return &multiplyExpression{
		binaryExpression: binaryExpression{
			lhs:          lhs,
			rhs:          rhs,
			isSimplified: true,
		},
	}
}

func (e *multiplyExpression) String() string {
	return fmt.Sprintf("(%s * %s)", e.lhs.String(), e.rhs.String())
}

func (e *namedExpression) Evaluate(inputs []int) int {
	return e.expr.Evaluate(inputs)
}

func (e *namedExpression) FindInputs(target int, d decider) (map[int]int, error) {
	return e.expr.FindInputs(target, d)
}

func (e *namedExpression) Range() IntRange {
	return e.expr.Range()
}

func (e *namedExpression) Simplify() Expression {
	return &namedExpression{
		name: e.name,
		expr: e.expr.Simplify(),
	}
}

func (e *namedExpression) String() string {
	return fmt.Sprintf("<%s>", e.name)
}

////////////////////////////////////////////////////////////////////////////////
// registers

func (r *registers) set(name string, value Expression) {
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

func (r *registers) get(name string) Expression {
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

////////////////////////////////////////////////////////////////////////////////
// parseInput

func parseInput(r io.Reader) *registers {

	zero := literalExpression{0}

	result := registers{
		w: &zero,
		x: &zero,
		y: &zero,
		z: &zero,
	}

	inputIndex := 0

	s := bufio.NewScanner(r)
	lineIndex := 0

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		lineIndex++

		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, " ")

		if len(parts) == 2 {

			// this must be an "inp"
			if parts[0] != "inp" {
				panic(fmt.Sprintf("Invalid unary operation: %s", parts[0]))
			}
			expr := inputExpression{inputIndex}

			inputIndex++

			result.set(parts[1], &expr)

			continue
		}

		lhs := result.get(parts[1])

		var rhs Expression

		literalValue, err := strconv.ParseInt(parts[2], 10, 32)
		if err == nil {
			rhs = &literalExpression{
				value: int(literalValue),
			}
		} else {
			rhs = result.get(parts[2])
		}

		expr := makeBinaryExpression(parts[0], lhs, rhs)

		// set the value of the specified register to the expression
		result.set(parts[1], expr.Simplify())

	}

	return &result
}

func makeBinaryExpression(kind string, lhs Expression, rhs Expression) Expression {
	switch kind {
	case "add":
		return &addExpression{
			binaryExpression: binaryExpression{
				lhs: lhs,
				rhs: rhs,
			},
		}
	case "div":
		return &divideExpression{
			binaryExpression: binaryExpression{
				lhs: lhs,
				rhs: rhs,
			},
		}

	case "eql":
		return &equalsExpression{
			binaryExpression: binaryExpression{
				lhs: lhs,
				rhs: rhs,
			},
		}

	case "mod":
		return &moduloExpression{
			binaryExpression: binaryExpression{
				lhs: lhs,
				rhs: rhs,
			},
		}

	case "mul":
		return &multiplyExpression{
			binaryExpression: binaryExpression{
				lhs: lhs,
				rhs: rhs,
			},
		}

	default:
		panic(fmt.Sprintf("Invalid op: %s", kind))
	}
}

// Take a binary expression (e.g. a +, *, /, etc.) and find inputs required
// to get it to equal <target>
func findInputsForBinaryExpression(
	e *binaryExpression,
	target int,
	getRhsValues func(lhsValue int, rhsRange IntRange) ([]int, error),
	d decider,
) (map[int]int, error) {

	lhsRange := e.lhs.Range()
	rhsRange := e.rhs.Range()

	var best map[int]int

	// for each value in left side's range, look for a corresponding value in the
	// right side's range and figure out the inputs needed to get them both to go there
	for lhsValue := lhsRange.min; lhsValue <= lhsRange.max; lhsValue += lhsRange.step {
		potentialRhsValues, err := getRhsValues(lhsValue, rhsRange)

		fmt.Printf("%d - %v\n", lhsValue, potentialRhsValues)

		if err != nil {
			continue
		}

		for _, rhsValue := range potentialRhsValues {

			lhsInputs, err := e.lhs.FindInputs(lhsValue, d)

			if err != nil {
				continue
			}

			rhsInputs, err := e.rhs.FindInputs(rhsValue, d)

			if err != nil {
				continue
			}

			bothSidesInSync := true
			inputs := make(map[int]int, len(lhsInputs)+len(rhsInputs))

			for index, value := range rhsInputs {
				lhsInputValue, lhsUsesInput := (lhsInputs)[index]
				if lhsUsesInput && lhsInputValue != value {
					// for this to work, left and right side need the same input set to
					// different values
					bothSidesInSync = false
					break
				}
				inputs[index] = value
			}

			if !bothSidesInSync {
				continue
			}

			for index, value := range lhsInputs {
				inputs[index] = value
			}

			if best == nil {
				best = inputs
			} else {
				b, err := d(best, inputs)
				if err == nil {
					best = b
				}
			}

		}
	}

	if best == nil {
		return nil, fmt.Errorf("No inputs can reach %d for ranges %v and %v", target, lhsRange, rhsRange)
	}

	return best, nil
}
