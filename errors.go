package dockerlang

import "errors"

var (
	TrivialWhitespaceError = errors.New("encountered a non-linebreak whitespace")
)
