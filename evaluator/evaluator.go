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

	// If the node is a *ast.Program, we evaluate the statements
	case *ast.Program:
		return evalProgram(node.Statements)

	/*
	 * Statements
	 */

	// If the node is a *ast.ExpressionStatement,
	// we evaluate the expression by recursively calling Eval()
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// If the node is a *ast.BlockStatement,
	// we evaluate the statements.
	// example: { <statement1>; <statement2>; ... }
	// or if-else blocks
	case *ast.BlockStatement:
		return evalBlockStatement(node)

	// If the node is a *ast.ReturnStatement,
	// we evaluate the return value
	// and return an object.ReturnValue object
	// which holds the value of the return statement
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: val}

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

	// If the node is a *ast.IfExpression,
	// we evaluate the condition and return the corresponding
	// consequence or alternative
	case *ast.IfExpression:
		return evalIfExpression(node)

	}

	// If we don't recognize the node, we return nil
	return nil
}

// evalProgram evaluates a slice of statements,
// returning the Eval() result of the last statement
func evalProgram(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)

		// If the result is a object.ReturnValue,
		// we break the loop and
		// return the unwraped value of the object.ReturnValue
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

// evalBlockStatement evaluates a block statement
// by evaluating each statement in the block
// and returning the Eval() result of the last statement
func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object
	for _, stmt := range block.Statements {
		result = Eval(stmt)

		// If the result is a object.ReturnValue,
		// we break the loop and
		// return the value of the object.ReturnValue
		// this is the wrapped value, must let evalProgram unwrap it
		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
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

// evaluates the infix expression for booleans
// by checking the operator returning the result.
// Example: <leftValue> <operator> <rightValue>
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

// evaluates the if expression by checking the condition
// and returning the consequence or alternative
// based on if the condition is true or not
// Example: if (<condition>) { <consequence> } else { <alternative> }
func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

// isTruthy checks if the object is truthy
// by checking if it is NULL, TRUE or FALSE
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
