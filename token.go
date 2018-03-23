package dockerlang

const (
	OPERATOR    = "OPERATOR"
	VARIABLE    = "VARIABLE"
	INT         = "INTEGER"
	PUNCTUATION = "PUNCTUATION" // parens
)

type Token struct {
	Type  string
	Value string
}
