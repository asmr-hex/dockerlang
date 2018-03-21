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
		Name: "src.doc",
		Body: &Expr{
			Op:    DIVISION_OPERATOR,
			Arity: OP_TO_ARITY[DIVISION_OPERATOR],
			Operands: []interface{}{
				&Expr{
					Op:    ADDITION_OPERATOR,
					Arity: OP_TO_ARITY[ADDITION_OPERATOR],
					Operands: []interface{}{&Expr{
						Op:       NOOP,
						Arity:    OP_TO_ARITY[NOOP],
						Operands: []interface{}{"2"},
					},
						&Expr{
							Op:       NOOP,
							Arity:    OP_TO_ARITY[NOOP],
							Operands: []interface{}{"3"},
						},
					},
				},
				&Expr{
					Op:       NOOP,
					Arity:    OP_TO_ARITY[NOOP],
					Operands: []interface{}{"1"},
				},
			},
		},
	}

	assert.EqualValues(t, expectedStackTree, compt.StackTree)
}

func TestParser_ImbalancedParens(t *testing.T) {
	compt := &Compterpreter{Config: &Config{SrcFileName: "src.doc"}}
	compt.Tokens = []Token{
		{Type: OPERATOR, Value: "+"},
		// is whint? ;)
		{Type: INT, Value: "2"},
		{Type: INT, Value: "3"},
		{Type: PUNCTUATION, Value: ")"},
	}

	err := compt.Parse()
	assert.EqualValues(t, err, DockerlangSyntaxError)
}

func TestParser_SyntaxError(t *testing.T) {
	compt := &Compterpreter{Config: &Config{SrcFileName: "src.doc"}}
	compt.Tokens = []Token{
		{Type: PUNCTUATION, Value: "("},
		{Type: OPERATOR, Value: "+"},
		{Type: OPERATOR, Value: "+"},
		{Type: INT, Value: "3"},
		{Type: PUNCTUATION, Value: ")"},
	}

	err := compt.Parse()
	assert.EqualValues(t, err, DockerlangSyntaxError)

	compt.Tokens = []Token{
		{Type: PUNCTUATION, Value: "("},
		{Type: INT, Value: "3"},
		{Type: OPERATOR, Value: "+"},
		{Type: PUNCTUATION, Value: ")"},
	}

	err = compt.Parse()
	assert.EqualValues(t, err, DockerlangSyntaxError)

	compt.Tokens = []Token{
		{Type: PUNCTUATION, Value: "("},
		{Type: OPERATOR, Value: "+"},
		{Type: INT, Value: "3"},
		{Type: INT, Value: "3"},
		{Type: INT, Value: "3"},
		{Type: INT, Value: "3"},
		{Type: INT, Value: "3"},
		{Type: INT, Value: "3"},
		{Type: INT, Value: "3"},
		{Type: INT, Value: "3"},
		{Type: INT, Value: "3"},
		{Type: PUNCTUATION, Value: ")"},
	}

	err = compt.Parse()
	assert.EqualValues(t, err, DockerlangSyntaxError)
}
