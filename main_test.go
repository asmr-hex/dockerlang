package dockerlang

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	retcode := m.Run()

	os.Exit(retcode)
}
