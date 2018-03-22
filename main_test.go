package dockerlang

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestMain(m *testing.M) {
	retcode := m.Run()

	os.Exit(retcode)
}

func TestCompterpreterSuite(t *testing.T) {
	suite.Run(t, new(CompterpreterSuite))
}

func TestLexerSuite(t *testing.T) {
	suite.Run(t, new(LexerSuite))
}

func TestParserSuite(t *testing.T) {
	suite.Run(t, new(ParserSuite))
}

func TestExecutionSuite(t *testing.T) {
	suite.Run(t, new(ExecutionSuite))
}
