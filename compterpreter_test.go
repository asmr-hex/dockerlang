package dockerlang

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	retcode := m.Run()

	os.Exit(retcode)
}

func TestLoadSourceCode_NoSuchFile(t *testing.T) {
	conf := &Config{SrcFileName: "nonexistent_test_src.doc"}
	compt := &Compterpreter{Config: conf}

	err := compt.LoadSourceCode()
	if err == nil {
		t.Error("failed to fail to find file")
	}
}

func TestLoadSourceCode(t *testing.T) {
	conf := &Config{SrcFileName: "test/test.doc"}
	compt := &Compterpreter{Config: conf}

	err := compt.LoadSourceCode()
	if err != nil {
		t.Error(err)
	}
}

func TestTokenizeNumber(t *testing.T) {
	conf := &Config{SrcFileName: "test/test.doc"}
	compt := &Compterpreter{Config: conf}

	err := compt.LoadSourceCode()
	if err != nil {
		t.Error(err)
	}

	// advance ptr to first character
	compt.Advance()

	compt.TokenizeNumber(compt.CurrentChar)

	if compt.CurrentToken != "1234" {
		t.Error("incorrect token!")
	}
}

func TestGetNextToken(t *testing.T) {
	conf := &Config{SrcFileName: "test/test.doc"}
	compt := &Compterpreter{Config: conf}

	err := compt.LoadSourceCode()
	if err != nil {
		t.Error(err)
	}

	// advance ptr to first character
	compt.Advance()

	s, err := compt.GetNextToken()
	if err != nil {
		t.Error(err)
	}
	if s != "1234" {
		t.Errorf("incorrect first token! Expected '1234' got '%s'", s)
	}

	s, err = compt.GetNextToken()
	if err != nil {
		t.Error(err)
	}
	if s != "5678" {
		t.Errorf("incorrect second token! Expected '5678' got '%s'", s)
	}
}
