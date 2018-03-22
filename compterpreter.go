package dockerlang

import (
	"bufio"
	"os"
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
	CurrentToken Token
	Symbols      *Symbols
	Tokens       []Token
	StackTree    *Expr
}

func NewCompterpreter(c *Config) *Compterpreter {
	// whenever we create a new compterpreter, we will also
	// create a new execution engine which is set in the global
	// scope.
	err := NewExecutionEngine()
	if err != nil {
		panic(err)
	}

	return &Compterpreter{
		Config:  c,
		Symbols: PopulateSymbols(),
	}
}

func (c *Compterpreter) Compterpret() error {
	var (
		err error
	)

	// always shutdown the docker execution engine
	defer ShutdownExecutionEngine()

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
	// Identifies tokens in the provided .doc code
	c.Lex()
	// Creates c.StackTree representing the provided .doc code
	c.Parse()
	// Actually dockerize and evaluate the StackTree
	c.Evaluate()

	return nil
}
