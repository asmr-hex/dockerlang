package dockerlang

import (
	"os"
)

type Config struct {
	ShowUsage   bool
	SrcFileName string
	BinFileName string
}

type Compiler struct {
	Config *Config
}

func Compile(c *Config) error {
	var (
		err error
	)

	compiler := &Compiler{
		Config: c,
	}

	err = compiler.ReadSource()
	if err != nil {
		return err
	}

	return nil
}

func (c *Compiler) ReadSource() error {
	// check to see if provided file exists
	info, err := os.Stat(c.Config.SrcFileName)
	if err != nil {
		return err
	}

	// TODO check filesize and permissions of file
	_ = info

	// stream file in loop
	// bufio.

	//

	return nil
}
