package dockerlang

import (
	"github.com/stretchr/testify/suite"
)

type CompterpreterSuite struct {
	suite.Suite
}

func (s *CompterpreterSuite) AfterTest(suiteName, testName string) {
	ShutdownExecutionEngine()
}

func (s *CompterpreterSuite) TestLoadSourceCode_NoSuchFile() {
	conf := &Config{SrcFileName: "nonexistent_test_src.doc"}
	compt := NewCompterpreter(conf)

	err := compt.LoadSourceCode()
	s.Error(err)
}

func (s *CompterpreterSuite) TestLoadSourceCode() {
	conf := &Config{SrcFileName: "test/test.doc"}
	compt := NewCompterpreter(conf)

	err := compt.LoadSourceCode()
	s.NoError(err)
}
