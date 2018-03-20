package dockerlang

const (
	OPERATOR    = "OPERATOR"
	INT         = "INTEGER"
	PUNCTUATION = "PUNCTUATION" // parens
)

type Token struct {
	Type  string
	Value string
}
