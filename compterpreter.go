package dockerlang

import (
	"bufio"
	"fmt"
	"os"
	"text/scanner"
	"unicode"
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

	compterpreter := &Compterpreter{
		Config: c,
	}

	// initialize a scanner to read through source code character
	// by character
	err = compterpreter.LoadSourceCode()
	if err != nil {
		return err
	}

	// start interpreting
	err = compterpreter.Interpret()
	if err != nil {
		return err
	}

	return nil
}

func (c *Compterpreter) LoadSourceCode() error {
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

	// set source code scanner
	reader := bufio.NewReader(fd)
	c.Scanner.Init(reader)

	return nil
}

func (c *Compterpreter) Interpret() error {
	// starting off from the beginning of the source file we will
	// always first advance to the first character.
	c.Advance()

	return nil
}

func (c *Compterpreter) GetNextToken() (string, error) {
	// we must clear the CurrentToken each time we get the next token!
	c.CurrentToken = ""

	// we are looping since there are characters we may want to ignore
	// for example, whitespace or something.
	for {
		switch {
		case c.IsWhitespace(c.CurrentChar):
			// igfnore whitespace
			c.Advance()
			continue
		// case c.IsOperator(c.CurrentChar):
		// 	// something
		// 	c.TokenizeOperator()
		// 	done = true
		case c.IsNumber(c.CurrentChar):
			// get full multidigit number token
			c.TokenizeNumber(c.CurrentChar)
		default:
			// we've encountered something very unexpected!
			// i'd like to panic, but i'm gunna keep my kewl
			return "", fmt.Errorf(
				"sry, but ive NO IDEA wut this char is: %s. Try typing another one(?)",
				string(c.CurrentChar),
			)
		}

		// if we ever get here, we have gotten the next token
		// and the CurrentChar should be pointing to the next
		// character we want to start tokenizing!
		break
	}

	// currently just returning the CurrentToken (what we just tokenized) and
	// assuming that the caller is appending the token to an array of seen tokens

	return c.CurrentToken, nil
}

func (c *Compterpreter) IsWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

func (c *Compterpreter) IsNumber(r rune) bool {
	return unicode.IsDigit(r)
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
