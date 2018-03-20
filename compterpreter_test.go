package dockerlang

import (
	"testing"
)

func TestLoadSourceCode_NoSuchFile(t *testing.T) {
	conf := &Config{SrcFileName: "nonexistent_test_src.doc"}
	compt := NewCompterpreter(conf)

	err := compt.LoadSourceCode()
	if err == nil {
		t.Error("failed to fail to find file")
	}
}

func TestLoadSourceCode(t *testing.T) {
	conf := &Config{SrcFileName: "test/test.doc"}
	compt := NewCompterpreter(conf)

	err := compt.LoadSourceCode()
	if err != nil {
		t.Error(err)
	}
}
