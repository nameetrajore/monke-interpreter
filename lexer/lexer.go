package lexer

import "monke/token"
// Lexer struct contains the string it is lexing
// the (current) 'position' it is at
// the (next position) 'readPosition'
// and the character at 'position'
type Lexer struct {
    input string //contains the entire program
    position int
    readPosition int
    ch byte
}

// helper function to skip all unnecessary white spaces
func (l* Lexer) skipWhitespace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r'{
        l.readChar()
    }
}

// returns the character at readPosition.
// If readPosition is beyond EOF it returns 0 (ASCII for EOF)
func (l* Lexer)peekChar() byte{
	if l.readPosition > len(l.input){
		return 0
	} else {
		return l.input[l.readPosition]
	}

}

// helper function to check if 'ch' is a digit
func isDigit(ch byte) bool{
	return '0' <= ch && ch <= '9'
}

// it returns the number it is reading and advances the pointers by calling readChar()
func (l* Lexer)readNumber() string {

	position := l.position
	for isDigit(l.ch){
		l.readChar()
	}

	return l.input[position:l.position]
}

// lexes the current 'ch' and creates new Tokens accordingly
func (l *Lexer) NextToken() token.Token {
    var tok token.Token

    l.skipWhitespace()

    switch l.ch {
	case '=':
		if l.peekChar() == '='{
		ch := l.ch
		l.readChar()
		//	for EQUALS
		tok.Literal = string(ch) + string(l.ch)
		tok.Type = token.EQ
		} else {
		// for ASSIGN
			tok = newToken(token.ASSIGN, l.ch)
		}
	case ';': tok = newToken(token.SEMICOLON, l.ch)
	case ':' : tok = newToken(token.COLON, l.ch)
	case '(': tok = newToken(token.LPAREN, l.ch)
	case ')': tok = newToken(token.RPAREN, l.ch)
	case '[': tok = newToken(token.LBRACKET, l.ch)
	case ']': tok = newToken(token.RBRACKET, l.ch)
	case ',': tok = newToken(token.COMMA, l.ch)
	case '+': tok = newToken(token.PLUS, l.ch)
	case '-': tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '='{
		ch := l.ch
		l.readChar()
			// for NOT EQUALS
		tok.Literal = string(ch) + string(l.ch)
		tok.Type = token.NOT_EQ
		} else {
			// for BANG
			tok = newToken(token.BANG, l.ch)
		}
	case '*': tok = newToken(token.ASTERISK, l.ch)
	case '/': tok = newToken(token.SLASH, l.ch)
	case '<': tok = newToken(token.LT, l.ch)
	case '>': tok = newToken(token.GT, l.ch)
	case '{': tok = newToken(token.LBRACE, l.ch)
	case '}': tok = newToken(token.RBRACE, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch){
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
    }
    l.readChar()
    return tok
}

// if nextToken is an identifier, readIdentifier is called to take the series of characters and return it
func (l *Lexer) readIdentifier() string {
    position := l.position
    for isLetter(l.ch){
        l.readChar()
    }

    return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1

	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break;
		}
	}

	return l.input[position : l.position]
}

// helper function to check if 'ch' is a letter
func isLetter (ch byte) bool {
    return 'a' <=ch && ch<= 'z' || 'A'<= ch && ch <='Z' || ch=='_'
}

// creates a new token
func newToken(tokenType token.TokenType, ch byte) token.Token {
    return token.Token{Type: tokenType, Literal: string(ch)}
}

// helper function to read the current character and advance readPosition and position
func (l* Lexer) readChar(){
	// sets current character to character at position
    if l.readPosition >= len(l.input){
        l.ch = 0
    } else {
        l.ch = l.input[l.readPosition]
    }
	// moves position by 1
    l.position = l.readPosition
	// moves readPosition by 1
    l.readPosition += 1
}

// creates a new lexer and returns a pointer to it
func New(input string) *Lexer{
    l := &Lexer{input: input}
    l.readChar() //XXX: Why?
    return l
}
