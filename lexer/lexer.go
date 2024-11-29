package lexer

import token "github.com/kev/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// gives us the next character and
// advances our position in the input string
func (l *Lexer) readChar() {

	//check if we've reached the end of the input
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	//move the position and advance readPosition up one
	l.position = l.readPosition
	l.readPosition += 1
}

// returns the next token in a token struct
// and advances the lexer's position in the input string
func (l *Lexer) NextToken() token.Token {
	var next_token token.Token

	l.eatWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			next_token = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			next_token = newToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			next_token = token.Token{Type: token.NE, Literal: string(ch) + string(l.ch)}
		} else {
			next_token = newToken(token.BANG, l.ch)
		}
	case ';':
		next_token = newToken(token.SEMICOLON, l.ch)
	case '(':
		next_token = newToken(token.LPAREN, l.ch)
	case ')':
		next_token = newToken(token.RPAREN, l.ch)
	case ',':
		next_token = newToken(token.COMMA, l.ch)
	case '+':
		next_token = newToken(token.PLUS, l.ch)
	case '{':
		next_token = newToken(token.LBRACE, l.ch)
	case '}':
		next_token = newToken(token.RBRACE, l.ch)
	case '-':
		next_token = newToken(token.MINUS, l.ch)
	case '*':
		next_token = newToken(token.ASTERISK, l.ch)
	case '/':
		next_token = newToken(token.SLASH, l.ch)
	case '<':
		next_token = newToken(token.LT, l.ch)
	case '>':
		next_token = newToken(token.GT, l.ch)
	case '"':
		next_token.Type = token.STRING
		next_token.Literal = l.readString()
	case 0:
		next_token.Literal = ""
		next_token.Type = token.EOF
	default:
		if isLetter(l.ch) {
			next_token.Literal = l.readIdentifier()
			next_token.Type = token.LookupIdent(next_token.Literal)
			return next_token
		} else if isDigit(l.ch) {
			next_token.Type = token.INT
			next_token.Literal = l.readNumber()
			return next_token
		} else {
			next_token = newToken(token.ILLEGAL, l.ch)
		}
	}
	l.readChar()
	return next_token
}

// helper function to create a new token
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// helper function to read an identifier
// and advance the lexer's position in the input string
// until it encounters a non-letter character
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// helper function to read a number
// and advance the lexer's position in the input string
// until it encounters a non-digit character
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// helper function to read a string
// and advance the lexer's position in the input string
// until it encounters a closing double quote
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// helper function to check if a character is a letter
// (a-z, A-Z, or an underscore)
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

// helper function to check if a character is a digit
// (0-9)
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// helper function to eat whitespaces
func (l *Lexer) eatWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// helper function to peek at the next character
// without advancing the lexer's position in the input string
func (l *Lexer) peekChar() byte {
	if l.readPosition > len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
