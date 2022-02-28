package d24

import (
	"fmt"
	"strings"
)

func PrettyPrintExpression(expr Expression, indent string) string {
	return prettyPrintExpressionAtIndent(expr, indent, 0)

}

func prettyPrintExpressionAtIndent(expr Expression, indent string, level int) string {
	var fullIndent string
	for i := 0; i < level; i++ {
		fullIndent = fullIndent + indent
	}

	switch expr := expr.(type) {
	case *AddExpression, *DivideExpression, *ModuloExpression, *MultiplyExpression:
		if isSimpleBinaryExpression(expr) {
			return fmt.Sprintf("%s%s", fullIndent, expr.String())
		} else {
			b := expr.(BinaryExpression)
			return strings.Join(
				[]string{
					fullIndent + "(",
					prettyPrintExpressionAtIndent(b.Lhs(), indent, level+1),
					fullIndent + indent + b.Operator(),
					prettyPrintExpressionAtIndent(b.Rhs(), indent, level+1),
					fullIndent + ")",
				},
				"\n",
			)
		}
	case *EqualsExpression:
		if isSimpleEqualsExpression(expr) {
			return fmt.Sprintf("%s%s", fullIndent, expr.String())
		} else {
			return strings.Join(
				[]string{
					fullIndent + "(",
					prettyPrintExpressionAtIndent(expr.Lhs(), indent, level+1),
					fullIndent + indent + "==",
					prettyPrintExpressionAtIndent(expr.Rhs(), indent, level+1),
					fullIndent + indent + "? 1 : 0",
					fullIndent + ")",
				},
				"\n",
			)
		}
	case *InputExpression, *LiteralExpression:
		return fmt.Sprintf("%s%s", fullIndent, expr)
	default:
		panic(fmt.Sprintf("Unhandled type: %T", expr))
	}
}

func indentString(value string, indent string) string {
	lines := strings.Split(value, "\n")
	for i := range lines {
		lines[i] = fmt.Sprintf("%s%s", indent, lines[i])
	}
	return strings.Join(lines, "\n")
}

func isSimpleEqualsExpression(expr Expression) bool {
	e, isEquals := expr.(*EqualsExpression)
	if !isEquals {
		return false
	}

	lhsSimple := isInputOrLiteralExpression(e.Lhs()) || isSimpleBinaryExpression(e.Lhs())
	rhsSimple := isInputOrLiteralExpression(e.Rhs()) || isSimpleBinaryExpression(e.Rhs())

	return lhsSimple && rhsSimple
}

func isSimpleBinaryExpression(expr Expression) bool {
	b, isBinary := expr.(BinaryExpression)

	if !isBinary {
		return false
	}

	return isInputOrLiteralExpression(b.Lhs()) && isInputOrLiteralExpression(b.Rhs())

}

func isInputOrLiteralExpression(expr Expression) bool {
	switch expr.(type) {
	case *InputExpression, *LiteralExpression:
		return true
	default:
		return false
	}
}
