package dockerlang

import (
	"github.com/stretchr/testify/suite"
)

type ParserSuite struct {
	suite.Suite
}

func (s *ParserSuite) AfterTest(suiteName, testName string) {
	ShutdownExecutionEngine()
}

func (s *ParserSuite) TestParser() {
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
	s.NoError(err)

	expectedTree := &Expr{
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

	s.EqualValues(expectedTree, compt.StackTree)
}

func (s *ParserSuite) TestParser_SyntaxError() {
	compt := &Compterpreter{Config: &Config{SrcFileName: "src.doc"}}
	compt.Tokens = []Token{
		{Type: PUNCTUATION, Value: "("},
		{Type: OPERATOR, Value: "+"},
		{Type: OPERATOR, Value: "+"},
		{Type: INT, Value: "3"},
		{Type: PUNCTUATION, Value: ")"},
	}

	err := compt.Parse()
	s.IsType(DockerlangSyntaxError{}, err)

	compt.Tokens = []Token{
		{Type: PUNCTUATION, Value: "("},
		{Type: INT, Value: "3"},
		{Type: OPERATOR, Value: "+"},
		{Type: PUNCTUATION, Value: ")"},
	}

	err = compt.Parse()
	s.IsType(DockerlangSyntaxError{}, err)

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
	s.IsType(DockerlangSyntaxError{}, err)
}
