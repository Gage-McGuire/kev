package evaluator

import (
	"fmt"

	"github.com/Gage-McGuire/kev/ast"
	"github.com/Gage-McGuire/kev/object"
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
// into an object.Object. It also takes an
// object.Enviroment to keep track of the variables
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// If the node is a *ast.Program, we evaluate start
	// by evaluating the statements in the program with evalProgram()
	case *ast.Program:
		return evalProgram(node.Statements, env)

	/*
	 * Statements
	 */

	// If the node is a *ast.ExpressionStatement,
	// we evaluate the expression by recursively calling Eval()
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	// If the node is a *ast.BlockStatement,
	// we evaluate the statements.
	// example: { <statement1>; <statement2>; ... }
	// or if-else blocks
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	// If the node is a *ast.ReturnStatement,
	// we evaluate the return value
	// and return an object.ReturnValue object
	// which holds the value of the return statement
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	// If the node is a *ast.VarStatement,
	// we evaluate the value of the var statement
	// and store it in the enviroment
	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

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
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	// If the node is a *ast.InfixExpression,
	// we evaluate the left and right side of the expression
	// and pass them to evalInfixExpression
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	// If the node is a *ast.IfExpression,
	// we evaluate the condition and return the corresponding
	// consequence or alternative
	case *ast.IfExpression:
		return evalIfExpression(node, env)

	// If the node is a *ast.CallExpression,
	// we evaluate the function and return the result
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

	// If the node is a *ast.IndexExpression,
	// we evaluate the left and index
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	/*
	 * Identifiers
	 */

	// If the node is a *ast.Identifier,
	// we evaluate the identifier
	case *ast.Identifier:
		return evalIdentifier(node, env)

	/*
	 * Literals
	 */

	// If the node is a *ast.FunctionLiteral,
	// we return an object.Function
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

	// If the node is a *ast.StringLiteral,
	// we return an object.String
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	// If the node is a *ast.ArrayLiteral,
	// we evaluate the elements and return an object.Array
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	}

	// If we don't recognize the node, we return nil
	return nil
}

// evalProgram evaluates a slice of statements,
// returning the Eval() result of the last statement
func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt, env)
		switch result := result.(type) {

		// If the result is a object.ReturnValue,
		// we break the loop and
		// return the unwraped value of the object.ReturnValue
		case *object.ReturnValue:
			return result.Value

		// If the result is a object.Error,
		// we break the loop and return the object.Error
		case *object.Error:
			return result
		}
	}

	return result
}

// evalBlockStatement evaluates a block statement
// by evaluating each statement in the block
// and returning the Eval() result of the last statement
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		// If the result is not nil and
		// is a object.ReturnValue or object.Error,
		// we break the loop and return the wrapped result
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
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
		return newError("unknown operator: %s%s", operator, right.Type())
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
	// we return a newError with the unknown operator
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
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

	switch {
	// If the left and right objects are not the same type,
	// we return a newError with the type mismatch
	// or unknown operator
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())

	// If the left and right objects are strings,
	// we evaluate the infix expression by calling evalStringInfixExpression
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evaluates the infix expression for strings
// by checking the operator and returning the result.
// Example: <leftValue> <operator> <rightValue>
func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	switch operator {
	case "==":
		return nativeBoolToBooleanObject(left.(*object.String).Value == right.(*object.String).Value)
	case "!=":
		return nativeBoolToBooleanObject(left.(*object.String).Value != right.(*object.String).Value)
	case "+":
		leftValue := left.(*object.String).Value
		rightValue := right.(*object.String).Value
		return &object.String{Value: leftValue + rightValue}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evaluates the if expression by checking the condition
// and returning the consequence or alternative
// based on if the condition is true or not
// Example: if (<condition>) { <consequence> } else { <alternative> }
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

// evaluates the identifier by checking if the identifier
// exists in the enviroment and returning the value or
// checking if the identifier is a builtin function
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

// Iterates over a slice of ast.Expressions and evaluates them
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

// applyFunction checks if the function is a *object.Function
// and applies the function by extending the environment or
// if the function is a *object.Builtin, it applies the function
func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Func(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

// extendFunctionEnv creates a new enclosed environment with the
// outer environment being set in the new environment as well
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

// unwrapReturnValue unwraps the return value to ensure
// the return value does not bubble up the stack and stop the
// execution of the program
func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
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

// creates a new object.Error
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// checks if the object is an error object
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
