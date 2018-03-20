package dockerlang

const (
	ADDITION_OPERATOR       = "+"
	SUBTRACTION_OPERATOR    = "†"
	MULTIPLICATION_OPERATOR = "*"
	DIVISION_OPERATOR       = "‡"
	MODULO_OPERATOR         = "%"
)

// all the language-defined tokens in dockerlang
type Symbols struct {
	Operators []string
	Keywords  []string
}

func PopulateSymbols() *Symbols {
	return &Symbols{
		Operators: []string{
			ADDITION_OPERATOR,
			SUBTRACTION_OPERATOR,
			MULTIPLICATION_OPERATOR,
			DIVISION_OPERATOR,
			MODULO_OPERATOR,
		},
		Keywords: []string{},
	}
}
