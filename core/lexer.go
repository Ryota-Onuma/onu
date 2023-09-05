package core

type Lexer struct {
	input        []rune
	position     int  // 入力における現在の位置（現在の文字を指し示す）
	readPosition int  // これから読み込む位置（現在の文字の次）
	ch           rune // 現在検査中の文字
	statements   []Statement
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: []rune(input)}
	l.readChar()
	l.addSemicoron()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = EOF
	} else {
		l.ch = rune(l.input[l.readPosition])
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=', '+', '-', '*', '/', '%', '(', ')', ';':
		tok = newToken(int(l.ch), l.ch)
	case '足':
		if l.peekChar() == 'す' {
			l.readChar()
			tok = newToken(int('+'), l.ch)
		} else {
			tok = newToken(int(UNKNOWN), l.ch)
		}
	case '引':
		if l.peekChar() == 'く' {
			l.readChar()
			tok = newToken(int('-'), l.ch)
		} else {
			tok = newToken(int(UNKNOWN), l.ch)
		}
	case '掛':
		if l.peekChar() == 'け' && l.peekChar2() == 'る' {
			l.readChar()
			l.readChar()
			tok = newToken(int('*'), l.ch)
		} else {
			tok = newToken(int(UNKNOWN), l.ch)
		}
	case '割':
		if l.peekChar() == 'る' {
			l.readChar()
			tok = newToken(int('/'), l.ch)
		} else {
			tok = newToken(int(UNKNOWN), l.ch)
		}
	case '余':
		if l.peekChar() == 'り' {
			l.readChar()
			tok = newToken(int('%'), l.ch)
		} else {
			tok = newToken(int(UNKNOWN), l.ch)
		}
	default:
		if l.ch == '定' {
			if l.peekChar() == '義' {
				l.readChar()
				tok = newToken(VAR, l.ch)
				return tok
			}
		} else if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = l.lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = INT
			return tok
		} else if l.ch == EOF {
			tok = newToken(EOF, l.ch)
		} else {
			tok = newToken(UNKNOWN, l.ch)
		}
	}
	l.readChar()
	return tok
}

func newToken(tokenType int, ch rune) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return EOF
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) peekChar2() rune {
	if l.readPosition+1 >= len(l.input) {
		return EOF
	} else {
		return l.input[l.readPosition+1]
	}
}

var keywords = map[string]int{
	"var": VAR,
}

func (l *Lexer) lookupIdent(ident string) int {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// 末尾にセミコロンがなかったら追加する
// TODO: if文とかにも対応
func (l *Lexer) addSemicoron() {
	for i, s := range l.input {
		// 末尾にセミコロンがなかったら追加する
		if i == len(l.input)-1 && s != ';' {
			l.input = append(l.input, ';')
		}
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return string(l.input[position:l.position])
}
