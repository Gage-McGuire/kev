package evaluator

import (
	"github.com/kev/ast"
	"github.com/kev/object"
)

var (
	// NULL is a singleton object.Null
	// representing the null value
	NULL = &object.Null{}

	// TRUE is a singleton object.Boolean
	// representing the boolean true
	TRUE = &object.Boolean{Value: true}

	// FALSE is a singleton object.Boolean
	// representing the boolean false
	FALSE = &object.Boolean{Value: false}
)

// Eval takes an AST node and evaluates it
// into an object.Object
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	/*
	 * Statements
	 */

	// If the node is a *ast.Program, we evaluate the statements
	case *ast.Program:
		return evalStatements(node.Statements)

	// If the node is a *ast.ExpressionStatement,
	// we evaluate the expression by recursively calling Eval()
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	/*
	 * Expressions
	 */

	// If the node is a *ast.IntegerLiteral,
	// we return an object.Integer
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	// If the node is a *ast.Boolean,
	// we return an object.Boolean
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	// If the node is a *ast.PrefixExpression,
	// we evaluate the right side of the expression
	// and pass it to evalPrefixExpression
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	// If the node is a *ast.InfixExpression,
	// we evaluate the left and right side of the expression
	// and pass them to evalInfixExpression
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	}

	// If we don't recognize the node, we return nil
	return nil
}

// evalStatements evaluates a slice of statements,
// returning the Eval() result of the last statement
func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

// evalPrefixExpression evaluates a prefix expression
// by checking the operator and passing the right object
// to the corresponding eval function
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

// evalBangOperatorExpression evaluates the bang operator
// by checking the right object and returning the inverse
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// converts a native Go boolean to a object.Boolean.
// This helps with the singleton pattern
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// evaluates the minus prefix operator by checking the right object
// and returning the negative value
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	// If the right object is not an object.Integer,
	// we return a NULL object
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}

	// We get the value of the right object
	value := right.(*object.Integer).Value

	// We return the address of the object.Integer
	// we created, that contains the negative value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	// If the left and right objects are integers,
	// we evaluate the infix expression by calling evalIntegerInfixExpression
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return evalIntegerInfixExpression(operator, left, right)
	}

	// If the left and right objects are booleans,
	// we evaluate the infix expression by calling evalBooleanInfixExpression
	if left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ {
		return evalBooleanInfixExpression(operator, left, right)
	}

	// If the left and right objects are not the same type,
	// we return NULL
	return NULL
}

// evaluates the infix expression for integers
// by checking the operator returning the result.
// Example: <leftValue> <operator> <rightValue>
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return NULL
	}
}

func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return NULL
	}
}
