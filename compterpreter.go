package dockerlang

import (
	"bufio"
	"fmt"
	"io"
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
	Symbols      *Symbols
	Tokens       []string
}

func NewCompterpreter(c *Config) *Compterpreter {
	return &Compterpreter{
		Config:  c,
		Symbols: PopulateSymbols(),
	}
}

func (c *Compterpreter) Compterpret() error {
	var (
		err error
	)

	// initialize a scanner to read through source code character
	// by character
	err = c.LoadSourceCode()
	if err != nil {
		return err
	}

	// start interpreting
	err = c.Interpret()
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
	c.Tokenize()

	c.Parse()

	return nil
}

func (c *Compterpreter) Tokenize() error {
	// starting off from the beginning of the source file we will
	// always first advance to the first character.
	c.Advance()
	for {
		token, err := c.GetNextToken()
		// gracefully catch end of file
		switch err {
		case io.EOF:
			return nil
		default:
			return err
		}

		c.Tokens = append(c.Tokens, token)
	}

	return nil
}

func (c *Compterpreter) Parse() error {
	return nil
}

func (c *Compterpreter) GetNextToken() (string, error) {
	var err error

	// we must clear the CurrentToken each time we get the next token!
	c.CurrentToken = ""

	// we are looping since there are characters we may want to ignore
	// for example, whitespace or something.
	for {
		switch {
		case c.IsWhitespace(c.CurrentChar):
			// igfnore whitespace
			err = c.Advance()
			if err != nil {
				return c.CurrentToken, err
			}
			continue
		case c.IsOperator(c.CurrentChar):
			// something
			err = c.TokenizeOperator(c.CurrentChar)
		case c.IsNumber(c.CurrentChar):
			// get full multidigit number token
			err = c.TokenizeNumber(c.CurrentChar)
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

	return c.CurrentToken, err
}

func (c *Compterpreter) IsWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

func (c *Compterpreter) IsNumber(r rune) bool {
	return unicode.IsDigit(r)
}

func (c *Compterpreter) IsOperator(r rune) bool {
	for _, symbol := range c.Symbols.Operators {
		if string(r) == symbol {
			return true
		}
	}
	return false
}

func (c *Compterpreter) TokenizeNumber(r rune) error {
	c.CurrentToken = c.CurrentToken + string(r)

	// check to see if we need to include the next character in the
	// current token
	if err := c.Advance(); err != nil {
		return err
	}

	if c.IsNumber(c.CurrentChar) {
		c.TokenizeNumber(c.CurrentChar)
	}

	return nil
}

func (c *Compterpreter) TokenizeOperator(r rune) error {
	c.CurrentToken = c.CurrentToken + string(r)
	if err := c.Advance(); err != nil {
		return err
	}

	// TODO: for a bright future which containts multi-symbol operators
	//if c.IsOperator(c.CurrentChar) {
	//	// check if the proposed multi-symbol operator is valid
	//	// if it's not, it's two operators next to each other
	//	c.TokenizeOperator(c.CurrentChar)
	//}
	return nil
}

func (c *Compterpreter) Advance() error {
	c.CurrentChar = c.Scanner.Next()

	if string(c.CurrentChar) == "EOF" {
		return io.EOF
	}

	return nil
}
