package dockerlang

import (
	"errors"
	"strings"
)

var (
	TrivialWhitespaceError = errors.New("encountered a non-linebreak whitespace")
	UnbalancedParenError   = errors.New("unbalanced paren")
	//DockerlangSyntaxError  = errors.New("what that's not how you write dockerlang come on homie")
)

type DockerlangSyntaxError struct {
	s string
}

func (e DockerlangSyntaxError) Error() string {
	return e.s
}

func SyntaxError(msgs ...string) error {
	message := "Dockerlang syntex error "

	return DockerlangSyntaxError{message + strings.Join(msgs, " ")}
}
