package dockerlang

import "regexp"

const (
	ADDITION_OPERATOR       = "+"
	SUBTRACTION_OPERATOR    = "†"
	MULTIPLICATION_OPERATOR = "*"
	DIVISION_OPERATOR       = "‡"
	MODULO_OPERATOR         = "%"
	VARIABLE_INITIALIZATION = "≡"
	VARIABLE_ASSIGNMENT     = "="
	EXIT_OPERATOR           = "ꙮ"
	IMPLICIT_LIST_OPERATOR  = "…"
	NOOP                    = "NOOP"

	R_PAREN_PUNCTION    = "("
	L_PAREN_PUNCUTATION = ")"

	VARIABLE_IDENTIFIER = "VARIABLE_IDENTIFIER"
	FUNCTION_IDENTIFIER = "FUNCTION_IDENTIFIER"

	IF_KEYWORD       = "if"
	FUNCTION_KEYWORD = "❦"
)

var (
	OP_TO_ARITY = map[string]int{
		ADDITION_OPERATOR:       2,
		SUBTRACTION_OPERATOR:    2,
		MULTIPLICATION_OPERATOR: 2,
		DIVISION_OPERATOR:       2,
		MODULO_OPERATOR:         2,
		VARIABLE_INITIALIZATION: 2,
		VARIABLE_ASSIGNMENT:     2,
		EXIT_OPERATOR:           1,
		NOOP:                    1,
	}

	VALID_IDENTIFIER_FIRST_SYMBOL = regexp.MustCompile("[a-zA-Z_]")
	VALID_IDENTIFIER_SYMBOL       = regexp.MustCompile("[a-zA-Z_\\d]")
)

// all the language-defined tokens in dockerlang
type Symbols struct {
	Operators   []string
	Keywords    []string
	Punctuation []string
}

func PopulateSymbols() *Symbols {
	return &Symbols{
		Operators: []string{
			ADDITION_OPERATOR,
			SUBTRACTION_OPERATOR,
			MULTIPLICATION_OPERATOR,
			DIVISION_OPERATOR,
			MODULO_OPERATOR,
			VARIABLE_INITIALIZATION,
			VARIABLE_ASSIGNMENT,
			EXIT_OPERATOR,
			IMPLICIT_LIST_OPERATOR,
			NOOP,
		},
		Keywords: []string{
			IF_KEYWORD,
			FUNCTION_KEYWORD,
		},
		Punctuation: []string{
			R_PAREN_PUNCTION,
			L_PAREN_PUNCUTATION,
		},
	}
}
