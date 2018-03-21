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

	expectedTree := &Expr{
		Name: "src.doc",
		Operands: []AST{
			&Expr{
				Op:    DIVISION_OPERATOR,
				Arity: OP_TO_ARITY[DIVISION_OPERATOR],
				Operands: []AST{
					&Expr{
						Op:    ADDITION_OPERATOR,
						Arity: OP_TO_ARITY[ADDITION_OPERATOR],
						Operands: []AST{
							&Literal{Type: INT, Value: "2"},
							&Literal{Type: INT, Value: "3"},
						},
					},
					&Literal{Type: INT, Value: "1"},
				},
			},
		},
	}

	assert.EqualValues(t, expectedTree, compt.StackTree)
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
