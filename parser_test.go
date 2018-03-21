package dockerlang

import (
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
		Name: "src.go",
		AST: &Expr{
			Op:    DIVISION_OPERATOR,
			Arity: OP_TO_ARITY[DIVISION_OPERATOR],
			LOperand: &Expr{
				Op:    ADDITION_OPERATOR,
				Arity: OP_TO_ARITY[ADDITION_OPERATOR],
				LOperand: &Expr{
					Op:       NOOP,
					Arity:    OP_TO_ARITY[NOOP],
					LOperand: 2,
				},
				ROperand: &Expr{
					Op:       NOOP,
					Arity:    OP_TO_ARITY[NOOP],
					LOperand: 3,
				},
			},
			ROperand: &Expr{
				Op:       NOOP,
				Arity:    OP_TO_ARITY[NOOP],
				LOperand: 1,
			},
		},
	}

	if expectedStackTree != compt.StackTree {
		t.Error("expected StackTree not equal to actual StackTree")
	}
}
