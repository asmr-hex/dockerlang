package dockerlang

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	retcode := m.Run()

	os.Exit(retcode)
}

func TestReadSource_NoSuchFile(t *testing.T) {
	conf := &Config{SrcFileName: "nonexistent_test_src.doc"}
	compt := &Compterpreter{Config: conf}

	err := compt.ReadSource()
	if err == nil {
		t.Error("failed to fail to find file")
	}
}

func TestReadSource(t *testing.T) {
	conf := &Config{SrcFileName: "test/test.doc"}
	compt := &Compterpreter{Config: conf}

	err := compt.ReadSource()
	if err != nil {
		t.Error(err)
	}
}
