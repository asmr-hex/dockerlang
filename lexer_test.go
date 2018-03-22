package dockerlang

import (
	"fmt"

	"github.com/stretchr/testify/suite"
)

type LexerSuite struct {
	suite.Suite
}

func (s *LexerSuite) AfterTest(suiteName, testName string) {
	ShutdownExecutionEngine()
}

func (s *LexerSuite) TestTokenizeNumber() {
	conf := &Config{SrcFileName: "test/test.doc"}
	compt := NewCompterpreter(conf)

	err := compt.LoadSourceCode()
	s.NoError(err)

	// advance ptr to first character
	compt.Advance()

	compt.TokenizeNumber(compt.CurrentChar)

	s.EqualValues(compt.CurrentToken.Value, "1234")
}

func (s *LexerSuite) TestGetNextToken() {
	conf := &Config{SrcFileName: "test/test.doc"}
	compt := NewCompterpreter(conf)

	err := compt.LoadSourceCode()
	s.NoError(err)

	// advance ptr to first character
	compt.Advance()

	t, err := compt.GetNextToken()
	s.NoError(err)

	s.EqualValues(t.Value, "1234")

	t, err = compt.GetNextToken()
	s.NoError(err)
	s.EqualValues(t.Value, "5678")
}

func (s *LexerSuite) TestIsOperator() {
	conf := &Config{SrcFileName: "test/test.doc"}
	compt := NewCompterpreter(conf)
	for _, operator := range []rune{'+', '†', '*', '‡', '%'} {
		ok := compt.IsOperator(operator)
		s.True(ok)
	}
	for _, operator := range []rune{'q', '!', '❧', '0', ' '} {
		ok := compt.IsOperator(operator)
		s.False(ok)
	}
}

func (s *LexerSuite) TestIsPunctuation() {
	conf := &Config{SrcFileName: "test/test.doc"}
	compt := NewCompterpreter(conf)
	for _, operator := range []rune{'(', ')', '(', ')'} {
		ok := compt.IsPunctuation(operator)
		s.True(ok)
	}
	for _, operator := range []rune{'q', '!', '❧', '0', ' '} {
		ok := compt.IsPunctuation(operator)
		s.False(ok)
	}
}

func (s *LexerSuite) TestTokenizeOperator() {
	conf := &Config{SrcFileName: "test/test-operators.doc"}
	compt := NewCompterpreter(conf)

	err := compt.LoadSourceCode()
	s.NoError(err)

	compt.Advance()
	// advance ptr to first character
	for _, op := range []string{"‡", "*", "+", "%", "†"} {
		compt.CurrentToken = Token{}
		compt.TokenizeOperator(compt.CurrentChar)
		if string(compt.CurrentChar) == "EOF" {
			break
		}
		s.EqualValues(compt.CurrentToken.Value, op)
	}
}

func (s *LexerSuite) TestLex() {
	conf := &Config{SrcFileName: "test/test_tokenize.doc"}
	compt := NewCompterpreter(conf)

	err := compt.LoadSourceCode()
	s.NoError(err)

	err = compt.Lex()
	s.NoError(err)

	for _, v := range compt.Tokens {
		fmt.Println(v.Value)
	}

	expectedTokens := []string{
		"\n", "123", "†", "3", "*", "2", "‡", "45787894357893", "\n", "0", "+", "00", "+", "1", "\n",
	}
	s.EqualValues(len(expectedTokens), len(compt.Tokens))
	for idx, token := range expectedTokens {
		s.EqualValues(token, compt.Tokens[idx].Value)
	}
}
