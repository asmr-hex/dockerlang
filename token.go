package dockerlang

const (
	OPERATOR    = "OPERATOR"
	IDENTIFIER  = "IDENTIFIER"
	KEYWORD     = "KEYWORD"
	INT         = "INTEGER"
	PUNCTUATION = "PUNCTUATION" // parens
)

type Token struct {
	Type  string
	Value string
}
