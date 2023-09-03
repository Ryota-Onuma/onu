
%{
// ヘッダー部
package main

import (
    "fmt"
    "text/scanner"
    "os"
    "strings"
    "strconv"
    "bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

)

// ast
type Expression interface{}

type NumberExpression struct {
	Literal string
}

// 二項演算子
type BinaryOperatorExpression struct {
	Left     Expression
	Operator rune
	Right    Expression
}

%}

%union{
    token Token
    expr  Expression
}
// %typeは非終端記号(プログラム、式、文)、
%type<expr> program
%type<expr> expr
// %tokenは終端記号(数値、演算子、括弧)
%token<token> NUMBER

%left '+', '-'
%left '*', '/'

// 規則部
%%

program
    : expr
    {
        $$ = $1
        yylex.(*Lexer).result = $$
    }

expr
    : NUMBER
    {
        $$ = NumberExpression{Literal: $1.literal}
    }
    | expr '+' expr
    {
        $$ = BinaryOperatorExpression{Left: $1, Operator: '+', Right: $3}
    }
    | expr '-' expr
    {
        $$ = BinaryOperatorExpression{Left: $1, Operator: '-', Right: $3}
    }
    | expr '*' expr
    {
        $$ = BinaryOperatorExpression{Left: $1, Operator: '*', Right: $3}
    }
    | expr '/' expr
    {
        $$ = BinaryOperatorExpression{Left: $1, Operator: '/', Right: $3}
    }
%%

// ユーザー定義


type Token struct {
	token   int
	literal string
}

type Lexer struct {
	scanner.Scanner
	result Expression
}

func (l *Lexer) Lex(lval *yySymType) int {
	token := int(l.Scan())
	if token == scanner.Int {
		token = NUMBER
	}
	lval.token = Token{token: token, literal: l.TokenText()}
	return token
}

func (l *Lexer) Error(e string) {
	panic(e)
}

func REPLExecute(input string) int {
	l := new(Lexer)
	l.Init(strings.NewReader(input))
	yyParse(l)
	return Eval(l.result)
}


func Eval(e Expression) int {
	switch e.(type) {
	case BinaryOperatorExpression:
		left := Eval(e.(BinaryOperatorExpression).Left)
		right := Eval(e.(BinaryOperatorExpression).Right)

		switch e.(BinaryOperatorExpression).Operator {
		case '+':
			return left + right
		case '-':
			return left - right
		case '*':
			return left * right
		case '/':
			return left / right
		}
	case NumberExpression:
		num, _ := strconv.Atoi(e.(NumberExpression).Literal)
		return num
	}
	return 0
}


func main() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    // 対話型のRPELループを開始
    reader := bufio.NewReader(os.Stdin)
    done := make(chan struct{})

    go func() {
        for {
            fmt.Print(">> ")
            input, err := reader.ReadString('\n')
            if err != nil {
                done <- struct{}{}
                return
            }

            // 改行文字を削除
            input = strings.TrimSpace(input)

            // 入力が "exit" だったらループを終了
            if input == "exit" {
                fmt.Println("RPELを終了します。")
                done <- struct{}{}
                return
            }

            // 入力を評価し、結果を表示
            result := REPLExecute(input)
		    fmt.Println(result)
        }
    }()

    select {
    case <-sigChan:
        os.Exit(1)
    case <-done:
    }
}
