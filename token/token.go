package token
type TokenType string
type Token struct {
    Type TokenType
    Literal string
}

//defining token types
const (
    ILLEGAL = "ILLEGAL"
    EOF = "EOF"

    //  Identifiers and literals
    IDENT = "IDENT"
    INT = "INT"
	STRING = "STRING"

    ASSIGN = "="
    PLUS = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

    COMMA = ","
    SEMICOLON = ";"
	COLON = ":"

    LPAREN = "("
    RPAREN = ")"
    LBRACE = "{"
    RBRACE = "{"
	LBRACKET = "["
	RBRACKET = "]"

    //Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
)

// to define and find token types for the given token(identifier)
var keywords = map[string]TokenType{
    "fn": FUNCTION,
    "let": LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
    if tok, ok := keywords[ident]; ok {
        return tok
    }
    return IDENT
}
