package core

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

type Env map[string]int

func Eval(stmt Statement, env Env) (string, error) {
	switch s := stmt.(type) {
	case ExpressionStatement:
		v, err := evalExpression(s.Expr, env)
		if err != nil {
			return "", err
		}
		return strconv.Itoa(v), nil
	case VarStatement:
		v, err := evalExpression(s.Value, env)
		if err != nil {
			return "", err
		}
		env[s.Name.Token.Literal] = v
		return "", nil
	}
	return "", errors.New("??? 謎のstatementです")
}

func evalExpression(expr Expression, env Env) (int, error) {
	switch e := expr.(type) {
	case BinaryOperatorExpression:
		left, err := evalExpression(e.Left, env)
		if err != nil {
			return 0, err
		}
		right, err := evalExpression(e.Right, env)
		if err != nil {
			return 0, err
		}

		switch e.Operator {
		case '+':
			return left + right, nil
		case '-':
			return left - right, nil
		case '*':
			return left * right, nil
		case '/':
			return left / right, nil
		case '%':
			return left % right, nil
		}
	case IntegerExpression:
		num, err := strconv.Atoi(e.Token.Literal)
		if err != nil {
			return 0, err
		}
		return num, nil
	case IdentifierExpression:
		if v, ok := env[e.Token.Literal]; ok {
			return v, nil
		} else {
			return 0, fmt.Errorf("%vは定義されてないですよ", e.Token.Literal)
		}
	}
	return 0, errors.New("??? 謎のexpressionです")
}

func Execute(input string, env Env) {
	statements := parse(input)
	for _, statement := range statements {
		s, err := Eval(statement, env)
		if err != nil {
			log.Fatal(err)
		}
		if s != "" {
			fmt.Println(s)
		}
	}
}
