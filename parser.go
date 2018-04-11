package dockerlang

import "fmt"

type Stack struct {
	Elements []AST
}

func NewStack() *Stack {
	return &Stack{Elements: []AST{}}
}

func (s *Stack) Push(e AST) {
	s.Elements = append(s.Elements, e)
}

func (s *Stack) Pop() AST {
	e := s.Peek()
	if e == nil {
		return nil
	}

	s.Elements = s.Elements[:len(s.Elements)-1]

	return e
}

func (s *Stack) Peek() AST {
	if len(s.Elements) == 0 {
		return nil
	}

	return s.Elements[len(s.Elements)-1]
}

func (s *Stack) Length() int {
	return len(s.Elements)
}

func (c *Compterpreter) Parse() error {
	// build the global StackTree, for all expressions in the global scope as part of an implicit anonymous function
	var (
		opsStack   = NewStack()
		exprStack  = NewStack()
		parenCount = 0
	)
	// TODO make the root expression an EXIT_OPERATOR which will operate
	// on all program roots.
	c.StackTree = &Expr{}

	for _, token := range c.Tokens {
		if parenCount < 0 {
			fmt.Println("UNBALANCED PARENS")
			// TODO (cw,mr|4.11.2018) be more specific -____-
			return DockerlangSyntaxError
		}

		switch token.Type {
		case OPERATOR:
			opsStack.Push(&Expr{Op: token.Value, Arity: OP_TO_ARITY[token.Value]})
		case INT:
			exprStack.Push(&Literal{Type: INT, Value: token.Value})
		case IDENTIFIER:
			// is this a declaration or a reference?
			identifier, err := c.StackTree.ParseIdentifier(token, opsStack)
			if err != nil {
				return err
			}

			exprStack.Push(identifier)
		case PUNCTUATION:
			switch token.Value {
			case "(":
				// TODO: eventually check a puntaution stack for syntax checking
				parenCount++
			case ")":
				// shit gets real
				var opsExpr = opsStack.Pop().(*Expr)
				// pop a count of arity items off exprStack
				for i := 0; i < opsExpr.Arity; i++ {
					// make sure we're not popping nil into exprs
					if exprStack.Peek() == nil {
						return DockerlangSyntaxError
					}
					opsExpr.Operands = append([]AST{exprStack.Pop()}, opsExpr.Operands...)
				}
				// push modified ops expr onto the expr stack
				exprStack.Push(opsExpr)

				parenCount--
			default:
				// whatever
			}
		}
		// if there is only 1 element in the expressionStack, we have successfully parsed
		// a single expression
		if exprStack.Length() == 1 && parenCount == 0 {
			// there should be nothing on the operations stack!
			if opsStack.Peek() != nil {
				fmt.Println("error in loop")
				fmt.Println(opsStack.Peek())
				fmt.Println(exprStack.Peek())
				// oh noooo!
				return DockerlangSyntaxError
			}

			// add this expression to the sequential list of expressions in the
			// programs execution
			c.StackTree.Operands = append(
				c.StackTree.Operands,
				exprStack.Pop().(*Expr),
			)
		}
	}

	// there should be nothing on the expression stack or operation stack
	if exprStack.Length() != 0 || opsStack.Peek() != nil {
		fmt.Println("SOMETHING IS AWRY")
		fmt.Println(opsStack.Peek())
		fmt.Println(exprStack.Peek())
		// oh no!
		return DockerlangSyntaxError
	}

	return nil
}

// TODO (cw,mr|4.11.2018) refactor this so that there is better separation of concerns
// i.e. it *would* be useful to have a function for adding local, global, args, etc. to
// an expression, but we might not want to do all this parsing in that function.
// NOTE: once we implement functions, we are going to want to check globals and args!
func (e *Expr) ParseIdentifier(token Token, opsStack *Stack) (*Identifier, error) {
	var (
		isDefined       bool = false
		knownIdentifier *Identifier
	)

	// check all locals to see if we've already defined this identifier
	for name, ast := range e.Locals {
		if token.Value == name {
			// this means we have already defined this identifer
			isDefined = true
			knownIdentifier = ast.(*Identifier)
			break
		}
	}

	prev := opsStack.Peek().(*Expr)

	// this is an identifier reference
	if prev.Op != VARIABLE_INITIALIZATION && prev.Op != FUNCTION_KEYWORD {
		if !isDefined {
			// TODO (cw,mr|4.11.2018) make this error more informative
			fmt.Println("TRYING TO USE AN UNDEFINED THING")
			return nil, DockerlangSyntaxError
		}

		// we are assuming that if an identifier is defined, then it is also bounded (or whatever)
		return knownIdentifier, nil
	}

	// if we are here, this is an identifier definition

	// we are trying to re-define this identifier
	if isDefined {
		fmt.Println("")
		return nil, DockerlangSyntaxError
	}

	// actually define this identifier
	switch prev.Op {
	case VARIABLE_INITIALIZATION:
		knownIdentifier = &Identifier{Type: VARIABLE_IDENTIFIER, Name: token.Value, Bound: true}
	case FUNCTION_KEYWORD:
		knownIdentifier = &Identifier{Type: FUNCTION_IDENTIFIER, Name: token.Value, Bound: true}
	}

	// TODO (cw, mr|4.11.2018) maybe put this in the Expr constructor once one exists
	if e.Locals == nil {
		e.Locals = map[string]AST{}
	}

	// add this to the local scope of the current expression
	e.Locals[knownIdentifier.Name] = knownIdentifier

	return knownIdentifier, nil
}
