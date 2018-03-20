package dockerlang

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"text/scanner"
)

type Config struct {
	ShowUsage   bool
	SrcFileName string
	BinFileName string
}

// is this a compiler or an interpreter? who knows????
type Compterpreter struct {
	Config       *Config
	Scanner      scanner.Scanner
	CurrentChar  rune
	CurrentToken string
	Tokens       []string
}

func Compterpret(c *Config) error {
	var (
		err error
	)

	compiler := &Compterpreter{
		Config: c,
	}

	err = compiler.ReadSource()
	if err != nil {
		return err
	}

	return nil
}

func (c *Compterpreter) ReadSource() error {
	// check to see if provided file exists
	info, err := os.Stat(c.Config.SrcFileName)
	if err != nil {
		return err
	}

	// TODO check filesize and permissions of file
	_ = info

	// open file
	fd, err := os.Open(c.Config.SrcFileName)
	if err != nil {
		return err
	}

	// stream file in loop
	reader := bufio.NewReader(fd)
	c.Scanner.Init(reader)

	c.Advance()

	t := c.GetNextToken()

	fmt.Println(t)

	return nil
}

func (c *Compterpreter) GetNextToken() string {
	var (
		done = false
	)

	for !done {
		switch {
		// case c.IsWhitespace(c.CurrentChar):
		// 	// igfnore whitespace
		// 	continue
		// case c.IsOperator(c.CurrentChar):
		// 	// something
		// 	c.TokenizeOperator()
		// 	done = true
		case c.IsNumber(c.CurrentChar):
			// get full multidigit number token
			c.TokenizeNumber(c.CurrentChar)
			done = true
		default:
			done = true
		}
	}

	return c.CurrentToken
}

func (c *Compterpreter) IsNumber(r rune) bool {
	_, err := strconv.Atoi(string(r))
	return err == nil
}

func (c *Compterpreter) TokenizeNumber(r rune) {
	c.CurrentToken = c.CurrentToken + string(r)

	// check to see if we need to include the next character in the
	// current token
	c.Advance()
	if c.IsNumber(c.CurrentChar) {
		c.TokenizeNumber(c.CurrentChar)
	}
}

func (c *Compterpreter) Advance() {
	c.CurrentChar = c.Scanner.Next()
}
