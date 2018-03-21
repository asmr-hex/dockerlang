package dockerlang

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParser(t *testing.T) {
	compt := &Compterpreter{Config: &Config{SrcFileName: "src.doc"}}
	compt.Tokens = []Token{
		{Type: PUNCTUATION, Value: "("},
		{Type: OPERATOR, Value: "‡"},
		{Type: PUNCTUATION, Value: "("},
		{Type: OPERATOR, Value: "+"},
		// is whint? ;)
		{Type: INT, Value: "2"},
		{Type: INT, Value: "3"},
		{Type: PUNCTUATION, Value: ")"},
		{Type: INT, Value: "1"},
		{Type: PUNCTUATION, Value: ")"},
	}

	err := compt.Parse()
	if err != nil {
		t.Error(err)
	}

	expectedStackTree := &StackTree{
		Name:  "src.doc",
		Scope: NewScope(),
		AST: &Expr{
			Op:    DIVISION_OPERATOR,
			Arity: OP_TO_ARITY[DIVISION_OPERATOR],
			LOperand: &Expr{
				Op:    ADDITION_OPERATOR,
				Arity: OP_TO_ARITY[ADDITION_OPERATOR],
				LOperand: &Expr{
					Op:       NOOP,
					Arity:    OP_TO_ARITY[NOOP],
					ROperand: "2",
				},
				ROperand: &Expr{
					Op:       NOOP,
					Arity:    OP_TO_ARITY[NOOP],
					ROperand: "3",
				},
			},
			ROperand: &Expr{
				Op:       NOOP,
				Arity:    OP_TO_ARITY[NOOP],
				ROperand: "1",
			},
		},
	}

	assert.EqualValues(t, expectedStackTree, compt.StackTree)
}
