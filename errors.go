package dockerlang

import (
	"errors"
)

var (
	TrivialWhitespaceError = errors.New("encountered a non-linebreak whitespace")
	UnbalancedParenError   = errors.New("unbalanced paren")
	//DockerlangSyntaxError  = errors.New("what that's not how you write dockerlang come on homie")
)

type errorString struct {
	s string
}

func (e errorString) Error() string {
	return e.s
}

func DockerlangSyntaxError(text string) error {
	return errorString{text}
}
