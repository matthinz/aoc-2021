package d24

import (
	"fmt"
	"testing"
)

func TestFormatInputExpression(t *testing.T) {
	t.Skip()
	expr := NewInputExpression(5)
	expected := "i5"
	actual := PrettyPrintExpression(expr, "\t")
	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestFormatLiteralExpression(t *testing.T) {
	expr := NewLiteralExpression(5)
	expected := "5"
	actual := PrettyPrintExpression(expr, "\t")
	if actual != expected {
		fmt.Println(actual)
		fmt.Println(expected)
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestFormatAddExpressionWithLiteralAndInput(t *testing.T) {
	expr := NewAddExpression(5, NewInputExpression(4))
	expected := "(5 + i4)"
	actual := PrettyPrintExpression(expr, "\t")
	if actual != expected {
		fmt.Println(actual)
		fmt.Println(expected)
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestFormatAddExpressionWithSubexpressions(t *testing.T) {
	expr := NewAddExpression(NewMultiplyExpression(5, NewInputExpression(0)), NewDivideExpression(6, NewInputExpression(1)))
	expected := "(\n\t(5 * i0)\n\t+\n\t(6 / i1)\n)"
	actual := PrettyPrintExpression(expr, "\t")

	if actual != expected {

		fmt.Println(actual)
		fmt.Println(expected)

		t.Errorf("Expected %s, got %s", expected, actual)
	}
}
