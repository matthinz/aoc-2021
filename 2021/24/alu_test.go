package d24

// import (
// 	_ "embed"
// 	"math"
// 	"strings"
// 	"testing"
// )

// type findInputsTest struct {
// 	name        string
// 	lhs         Expression
// 	rhs         Expression
// 	target      int
// 	pref        decider
// 	expected    map[int]int
// 	expectError bool
// }

// //go:embed input
// var realInput string

// func TestRealInput(t *testing.T) {
// 	parseInput(strings.NewReader(realInput))
// }

// func TestAlu(t *testing.T) {
// 	input := `
// inp w
// add z w
// mod z 2
// div w 2
// add y w
// mod y 2
// div w 2
// add x w
// mod x 2
// div w 2
// mod w 2
// 	`
// 	r := parseInput(strings.NewReader(input))

// 	inputs := []int{10}

// 	w := r.w.Evaluate(inputs)
// 	x := r.x.Evaluate(inputs)
// 	y := r.y.Evaluate(inputs)
// 	z := r.z.Evaluate(inputs)

// 	if w != 1 {
// 		t.Fatalf("expected w to be 1, but was %d", w)
// 	}

// 	if x != 0 {
// 		t.Fatalf("expected x to be 0, but was %d", x)
// 	}

// 	if y != 1 {
// 		t.Fatalf("expected y to be 1, but was %d", y)
// 	}

// 	if z != 0 {
// 		t.Fatalf("expected z to be 0, but was %d", z)
// 	}

// }

// func TestAluEquality(t *testing.T) {
// 	input := `
// inp w
// add x w
// add x 5
// eql x 7
// 	`
// 	r := parseInput(strings.NewReader(input))

// 	inputs := []int{2}

// 	t.Logf("x: %s", r.x.String())

// 	w := r.w.Evaluate(inputs)
// 	x := r.x.Evaluate(inputs)
// 	y := r.y.Evaluate(inputs)
// 	z := r.z.Evaluate(inputs)

// 	if w != 2 {
// 		t.Fatalf("expected w to be 10, but was %d", w)
// 	}

// 	if x != 1 {
// 		t.Fatalf("expected x to be 1, but was %d", x)
// 	}

// 	if y != 0 {
// 		t.Fatalf("expected y to be 0, but was %d", y)
// 	}

// 	if z != 0 {
// 		t.Fatalf("expected z to be 0, but was %d", z)
// 	}
// }

// func TestSimplifyEquality(t *testing.T) {
// 	e := equalsExpression{
// 		binaryExpression: binaryExpression{
// 			lhs: &addExpression{
// 				binaryExpression: binaryExpression{
// 					lhs: &inputExpression{
// 						index: 0,
// 					},
// 					rhs: &literalExpression{
// 						value: 2,
// 					},
// 				},
// 			},
// 			rhs: &literalExpression{
// 				value: 7,
// 			},
// 		},
// 	}

// 	expected := "((i0 + 2) == 7 ? 1 : 0)"
// 	if e.String() != expected {
// 		t.Fatalf("Expression is wrong. Expected %s, but got %s", expected, e.String())
// 	}

// 	simplified := e.Simplify()
// 	if simplified.String() != e.String() {
// 		t.Fatalf("simplified expression was wrong. expected %s, got %s", e.String(), simplified.String())
// 	}

// }

// func TestSimplifyModuloWithLiterals(t *testing.T) {
// 	expr := moduloExpression{
// 		binaryExpression: binaryExpression{
// 			lhs: &literalExpression{
// 				value: 10,
// 			},
// 			rhs: &literalExpression{
// 				value: 3,
// 			},
// 		},
// 	}

// 	simplified := expr.Simplify()
// 	expected := "1"
// 	if simplified.String() != expected {
// 		t.Errorf("simplified wrong. expected %s, got %s", expected, simplified.String())
// 	}
// }

// func TestSimplifyModuloWithInput(t *testing.T) {
// 	expr := moduloExpression{
// 		binaryExpression: binaryExpression{
// 			lhs: &literalExpression{
// 				value: 10,
// 			},
// 			rhs: &inputExpression{
// 				index: 0,
// 			},
// 		},
// 	}

// 	r := expr.Range()

// 	expected := IntRange{0, 4}

// 	if r != expected {
// 		t.Errorf("Range was wrong. Expected %v, got %v", expected, r)
// 	}

// }

// func TestAdditionFindInputs(t *testing.T) {

// 	tests := []findInputsTest{
// 		{
// 			name:     "FindSingleRhsInput",
// 			lhs:      &literalExpression{5},
// 			rhs:      &inputExpression{0},
// 			target:   12,
// 			pref:     preferFirstInput,
// 			expected: map[int]int{0: 7},
// 		},
// 		{
// 			name:     "FindSingleLhsInput",
// 			lhs:      &inputExpression{0},
// 			rhs:      &literalExpression{5},
// 			target:   12,
// 			pref:     preferFirstInput,
// 			expected: map[int]int{0: 7},
// 		},
// 		{
// 			name:     "FindLhsAndRhsInputsPreferringFirst",
// 			lhs:      &inputExpression{0},
// 			rhs:      &inputExpression{1},
// 			target:   12,
// 			pref:     preferFirstInput,
// 			expected: map[int]int{0: 3, 1: 9},
// 		},
// 		{
// 			name:     "FindLhsAndRhsInputsPreferringLargerNumber",
// 			lhs:      &inputExpression{0},
// 			rhs:      &inputExpression{1},
// 			target:   12,
// 			pref:     preferInputsThatMakeLargerNumber,
// 			expected: map[int]int{0: 9, 1: 3},
// 		},
// 	}

// 	factory := func(lhs, rhs Expression) Expression {
// 		return &addExpression{
// 			binaryExpression: binaryExpression{
// 				lhs: lhs,
// 				rhs: rhs,
// 			},
// 		}

// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, makeFindInputsTest(test, factory))
// 	}
// }

// func TestMultiplicationFindInputs(t *testing.T) {
// 	tests := []findInputsTest{
// 		{
// 			name:     "FindSingleRhsInput",
// 			lhs:      &literalExpression{5},
// 			rhs:      &inputExpression{0},
// 			target:   40,
// 			pref:     preferFirstInput,
// 			expected: map[int]int{0: 8},
// 		},
// 		{
// 			name:     "FindSingleLhsInput",
// 			lhs:      &inputExpression{0},
// 			rhs:      &literalExpression{5},
// 			target:   40,
// 			pref:     preferFirstInput,
// 			expected: map[int]int{0: 8},
// 		},
// 		{
// 			name:     "FindLhsAndRhsInputsPreferringFirst",
// 			lhs:      &inputExpression{0},
// 			rhs:      &inputExpression{1},
// 			target:   24,
// 			pref:     preferFirstInput,
// 			expected: map[int]int{0: 3, 1: 8},
// 		},
// 		{
// 			name:     "FindLhsAndRhsInputsPreferringLargerNumber",
// 			lhs:      &inputExpression{0},
// 			rhs:      &inputExpression{1},
// 			target:   24,
// 			pref:     preferInputsThatMakeLargerNumber,
// 			expected: map[int]int{0: 8, 1: 3},
// 		},
// 		{
// 			name:     "FindRhsZero",
// 			lhs:      &literalExpression{0},
// 			rhs:      &inputExpression{0},
// 			target:   0,
// 			pref:     preferFirstInput,
// 			expected: map[int]int{0: 1},
// 		},
// 		{
// 			name:     "FindLhsZero",
// 			lhs:      &inputExpression{0},
// 			rhs:      &literalExpression{0},
// 			target:   0,
// 			pref:     preferFirstInput,
// 			expected: map[int]int{0: 1},
// 		},
// 		{
// 			name:        "FindBothZero",
// 			lhs:         &inputExpression{0},
// 			rhs:         &inputExpression{0},
// 			target:      0,
// 			pref:        preferFirstInput,
// 			expectError: true,
// 		},
// 		{
// 			name:     "FindLhs1",
// 			lhs:      &inputExpression{0},
// 			rhs:      &literalExpression{1},
// 			target:   9,
// 			pref:     preferFirstInput,
// 			expected: map[int]int{0: 9},
// 		},
// 		{
// 			name:     "FindRhs1",
// 			lhs:      &literalExpression{1},
// 			rhs:      &inputExpression{0},
// 			target:   9,
// 			pref:     preferFirstInput,
// 			expected: map[int]int{0: 9},
// 		},
// 	}

// 	factory := func(lhs, rhs Expression) Expression {
// 		return &multiplyExpression{
// 			binaryExpression: binaryExpression{
// 				lhs: lhs,
// 				rhs: rhs,
// 			},
// 		}
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, makeFindInputsTest(test, factory))
// 	}

// }

// func makeFindInputsTest(test findInputsTest, f func(lhs, rhs Expression) Expression) func(*testing.T) {
// 	return func(t *testing.T) {
// 		expr := f(test.lhs, test.rhs)

// 		inputs, err := expr.FindInputs(test.target, test.pref)

// 		if test.expectError {
// 			if err == nil {
// 				t.Fatal("Expected error, but didn't get one")
// 			}
// 			return
// 		}

// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		if test.expected == nil && inputs != nil {
// 			t.Errorf("Expected nil inputs but got %v", inputs)
// 			return
// 		}

// 		if inputs == nil {
// 			t.Errorf("Expected %v but got nil", test.expected)
// 			return
// 		}

// 		if len(inputs) != len(test.expected) {
// 			t.Errorf("Expected %v (%d), got %v (%d)", test.expected, len(test.expected), inputs, len(inputs))
// 			return
// 		}

// 		for inputIndex, expectedValue := range test.expected {
// 			actualValue, inputIsFound := inputs[inputIndex]
// 			if !inputIsFound || actualValue != expectedValue {
// 				t.Errorf("Expected %v (%d), got %v (%d)", test.expected, len(test.expected), inputs, len(inputs))
// 				return
// 			}
// 		}
// 	}
// }

// func preferFirstInput(a, b map[int]int) (map[int]int, error) {
// 	return a, nil
// }

// func preferSecondInput(a, b map[int]int) (map[int]int, error) {
// 	return b, nil
// }

// func preferInputsThatMakeLargerNumber(a, b map[int]int) (map[int]int, error) {

// 	aValue := inputsToNumber(a)
// 	bValue := inputsToNumber(b)

// 	if aValue >= bValue {
// 		return a, nil
// 	} else {
// 		return b, nil
// 	}
// }

// func inputsToNumber(inputs map[int]int) int {
// 	var digits []int

// 	for inputIndex, inputValue := range inputs {
// 		if len(digits) < inputIndex+1 {
// 			temp := make([]int, inputIndex+1)
// 			copy(temp, digits)
// 			digits = temp
// 		}
// 		digits[inputIndex] = inputValue
// 	}

// 	var result int

// 	for i, digit := range digits {
// 		power := (len(digits) - (i + 1))
// 		result += int(math.Pow10(power)) * digit
// 	}

// 	return result
// }
