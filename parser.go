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
		opsStack  = NewStack()
		exprStack = NewStack()
	)
	// TODO make the root expression an EXIT_OPERATOR which will operate
	// on all program roots.
	c.StackTree = &Expr{}

	for _, token := range c.Tokens {
		switch token.Type {
		case OPERATOR:
			opsStack.Push(&Expr{Op: token.Value, Arity: OP_TO_ARITY[token.Value]})
		case INT:
			exprStack.Push(&Literal{Type: INT, Value: token.Value})
		case PUNCTUATION:
			switch token.Value {
			case "(":
				// TODO: eventually check a puntaution stack for syntax checking
			case ")":
				// shit gets real
				var opsExpr = opsStack.Pop().(*Expr)
				// pop a count of arity items off exprStack
				for i := 0; i < opsExpr.Arity; i++ {
					// make sure we're not popping nil into exprs
					if exprStack.Peek() == nil {
						fmt.Println("1")
						return DockerlangSyntaxError
					}
					opsExpr.Operands = append([]AST{exprStack.Pop()}, opsExpr.Operands...)
				}
				// push modified ops expr onto the expr stack
				exprStack.Push(opsExpr)
			default:
				// whatever
			}
		}
	}

	// there should only be one expr in exprStack
	if exprStack.Length() != 1 {
		// oh no!
		return DockerlangSyntaxError
	}
	if opsStack.Peek() != nil {
		fmt.Println(opsStack.Peek())
		// oh noooo!
		return DockerlangSyntaxError
	}

	fmt.Println(exprStack)

	c.StackTree.Operands = []AST{exprStack.Pop().(*Expr)}

	return nil
}
