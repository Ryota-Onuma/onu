package core

type Node interface{}
type Expression interface {
	Node
	expression()
}

type IntegerExpression struct {
	Token Token
}

// 二項演算子
type BinaryOperatorExpression struct {
	Left     Expression
	Operator int
	Right    Expression
}

type IdentifierExpression struct {
	Token Token
}

func (n IntegerExpression) expression()        {}
func (b BinaryOperatorExpression) expression() {}
func (v IdentifierExpression) expression()     {}

type Statement interface {
	Node
	statement()
}

type ExpressionStatement struct {
	Expr Expression
}

type VarStatement struct {
	Name  IdentifierExpression
	Value Expression
}

func (e ExpressionStatement) statement() {}
func (v VarStatement) statement()        {}
