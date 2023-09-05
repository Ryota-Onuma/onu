
%{
// ヘッダー部
package core

import (
    "fmt"
)
%}

%union{
    tok Token
    expr  Expression
    statement  Statement
    statements []Statement
}
%type<statements> statements
%type<statement> statement
%type<expr> expr

%token<tok> IDENT INT VAR

%left '+' '-'
%left '*' '/' '%'

// 規則部
%%

statements
        :
        {
            $$ = nil
            if l, ok := yylex.(*LexerWrapper); ok {
                l.lexer.statements = $$
            }
        }
        | statement statements
        {
            $$ = append([]Statement{$1}, $2...)
            if l, ok := yylex.(*LexerWrapper); ok {
                l.lexer.statements = $$
            }
        }
statement
        : expr ';'
        {
            $$ = ExpressionStatement{Expr: $1}
        }
        | VAR IDENT '=' expr ';'
        {
            $$ = VarStatement{Name: IdentifierExpression{Token: $2}, Value: $4}
        }
expr    
        : INT
        {
            $$ = IntegerExpression{Token: $1}
        } | IDENT
        {
            $$ = IdentifierExpression{Token: $1}
        } | expr '+' expr
        {
            $$ = BinaryOperatorExpression{Left: $1, Operator: int('+'), Right: $3}
        } | expr '-' expr
        {
            $$ = BinaryOperatorExpression{Left: $1, Operator: int('-'), Right: $3}
        } | expr '*' expr
        {
            $$ = BinaryOperatorExpression{Left: $1, Operator: int('*'), Right: $3}
        } | expr '/' expr
        {
            $$ = BinaryOperatorExpression{Left: $1, Operator: int('/'), Right: $3}
        } | expr '%' expr
        {
            $$ = BinaryOperatorExpression{Left: $1, Operator: int('%'), Right: $3}
        }
%%

// ユーザー定義

type LexerWrapper struct {
    lexer *Lexer
}

func (lw *LexerWrapper) Lex(lval *yySymType) int {
    token := lw.lexer.NextToken()
    if token.Type == EOF {
        return 0
    }
    if token.Type == UNKNOWN {
        panic("unknown token")
    }
    lval.tok = token
    return token.Type
}

func (lw *LexerWrapper) Error(s string) {
    fmt.Println(s)
}

func parse(input string) []Statement {
	l := &LexerWrapper{lexer: NewLexer(input)}
    yyParse(l)
	return l.lexer.statements
}