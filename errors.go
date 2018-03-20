package dockerlang

import "errors"

var (
	TrivialWhitespaceError = errors.New("encountered a non-linebreak whitespace")
	UnbalancedParenError   = errors.New("unbalanced paren")
	DockerlangSyntaxError  = errors.New("what that's not how you write dockerlang come on homie")
)
