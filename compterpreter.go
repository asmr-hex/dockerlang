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
	CurrentToken Token
	Symbols      *Symbols
	Tokens       []Token
	StackTree    *StackTree
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
	c.Lex()

	c.Parse()

	return nil
}

func (c *Compterpreter) Lex() error {
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
	// build the global StackTree, for all expressions in the global scope as part of an implicit anonymous function
	var (
		opsStack  = []Expr{}
		exprStack = []Expr{}
	)
	c.StackTree = NewStackTree(c.Config.SrcFileName)

	for _, token := range c.Tokens {
		switch token.Type {
		case OPERATOR:
			opsStack = append(opsStack, Expr{Op: token.Value, Arity: OP_TO_ARITY[token.Value]})
		case INT:
			exprStack = append(exprStack, Expr{Op: NOOP, Arity: OP_TO_ARITY[NOOP], lOperand: token.Value})
		case PUNCTUATION:
			switch token.Value {
			case "(":
				opsStack = append(opsStack, Expr{Op: token.Value, Arity: 1, lOperand: token.Value})
			case ")":
				// shit gets real
				var opsExpr = opsStack[len(opsStack)-1]
				// pop a count of arity items off exprStack
				var exprs = exprStack[len(exprStack)-1-opsExpr.Arity:]
				// load popped exprs into the ops expr
				opsExpr.lOperand = exprs[0]
				if len(exprs) > 1 {
					opsExpr.rOperand = exprs[1]
				}
				// update the stacks
				opsStack = opsStack[len(opsStack)-1:]
				// the +1 here should be the open paren for this expression
				var betterBeAnOpenParen = opsStack[len(opsStack)-1]
				if betterBeAnOpenParen.Op != "(" {
					return UnbalancedParenError
				}
				exprStack = exprStack[len(exprStack)-(opsExpr.Arity+1):]
				// push modified ops expr onto the expr stack
				exprStack = append(exprStack, opsExpr)
			default:
				// whatever
			}
		}
		// there should only be one expr in exprStack
		if len(exprStack) != 1 {
			// oh no!
			return DockerlangSyntaxError
		}
		if len(opsStack) > 0 {
			// oh noooo!
			return DockerlangSyntaxError
		}
		c.StackTree.AST = &exprStack[0]
	}

	return nil
}

func (c *Compterpreter) GetNextToken() (Token, error) {
	var err error

	// we must clear the CurrentToken each time we get the next token!
	c.CurrentToken = Token{}

	// we are looping since there are characters we may want to ignore
	// for example, whitespace or something.
	for {
		switch {
		case c.IsWhitespace(c.CurrentChar):
			// igfnore non-linebreak whitespace
			err = c.TokenizeWhitespace(c.CurrentChar)
			switch err {
			case io.EOF:
				return c.CurrentToken, err
			case TrivialWhitespaceError:
				continue
			}
		case c.IsOperator(c.CurrentChar):
			// something
			err = c.TokenizeOperator(c.CurrentChar)
		case c.IsNumber(c.CurrentChar):
			// get full multidigit number token
			err = c.TokenizeNumber(c.CurrentChar)
		default:
			// we've encountered something very unexpected!
			// i'd like to panic, but i'm gunna keep my kewl
			return Token{}, fmt.Errorf(
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

func (c *Compterpreter) TokenizeWhitespace(r rune) error {
	if r == '\n' {
		c.CurrentToken.Value = string(r)
		c.CurrentToken.Type = PUNCTUATION
		return nil
	}

	if err := c.Advance(); err != nil {
		return err
	}

	return TrivialWhitespaceError
}

func (c *Compterpreter) TokenizeNumber(r rune) error {
	c.CurrentToken.Type = INT
	c.CurrentToken.Value = c.CurrentToken.Value + string(r)

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
	c.CurrentToken.Type = OPERATOR
	c.CurrentToken.Value = c.CurrentToken.Value + string(r)
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
