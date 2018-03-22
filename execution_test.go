package dockerlang

import (
	"github.com/stretchr/testify/suite"
)

type ExecutionSuite struct {
	suite.Suite
}

func (s *ExecutionSuite) AfterTest(suiteName, testName string) {
	ShutdownExecutionEngine()
}

// NOTE: for these tests to run, we need to ensure that docker is running on the
// host machine!
func (s *ExecutionSuite) TestNewExecutionEngine() {
	err := NewExecutionEngine()
	s.NoError(err)
}

func (s *ExecutionSuite) TestShutdownExecutionEngine() {
	err := NewExecutionEngine()
	s.NoError(err)

	err = ShutdownExecutionEngine()
	s.NoError(err)
}

func (s *ExecutionSuite) TestShutdownExecutionEngine_NonExistent() {
	executer.Network = "non-existent network"

	err := ShutdownExecutionEngine()
	s.NoError(err)
}
