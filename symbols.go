package dockerlang

const (
	ADDITION_OPERATOR       = "+"
	SUBTRACTION_OPERATOR    = "†"
	MULTIPLICATION_OPERATOR = "*"
	DIVISION_OPERATOR       = "‡"
	MODULO_OPERATOR         = "%"
	NOOP                    = "NOOP"
)

var (
	OP_TO_ARITY = map[string]int{
		ADDITION_OPERATOR:       2,
		SUBTRACTION_OPERATOR:    2,
		MULTIPLICATION_OPERATOR: 2,
		DIVISION_OPERATOR:       2,
		MODULO_OPERATOR:         2,
		NOOP:                    1,
	}
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
			NOOP,
		},
		Keywords: []string{},
	}
}
