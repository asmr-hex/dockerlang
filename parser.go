package dockerlang

type ExprStack struct {
	Elements []*Expr
}

func NewExprStack() *ExprStack {
	return &ExprStack{Elements: []*Expr{}}
}

func (s *ExprStack) Push(e *Expr) {
	s.Elements = append(s.Elements, e)
}

func (s *ExprStack) Pop() *Expr {
	e := s.Peek()
	if e == nil {
		return nil
	}

	s.Elements = s.Elements[:len(s.Elements)-1]

	return e
}

func (s *ExprStack) Peek() *Expr {
	if len(s.Elements) == 0 {
		return nil
	}

	return s.Elements[len(s.Elements)-1]
}

func (s *ExprStack) Length() int {
	return len(s.Elements)
}

func (c *Compterpreter) Parse() error {
	// build the global StackTree, for all expressions in the global scope as part of an implicit anonymous function
	var (
		opsStack  = NewExprStack()
		exprStack = NewExprStack()
	)
	c.StackTree = NewStackTree(c.Config.SrcFileName)

	for _, token := range c.Tokens {
		switch token.Type {
		case OPERATOR:
			opsStack.Push(&Expr{Op: token.Value, Arity: OP_TO_ARITY[token.Value]})
		case INT:
			exprStack.Push(&Expr{Op: NOOP, Arity: OP_TO_ARITY[NOOP], ROperand: token.Value})
		case PUNCTUATION:
			switch token.Value {
			case "(":
				opsStack.Push(&Expr{Op: token.Value, Arity: 1, ROperand: token.Value})
			case ")":
				// shit gets real
				var opsExpr = opsStack.Pop()
				// pop a count of arity items off exprStack
				exprs := []*Expr{}
				for i := 0; i < opsExpr.Arity; i++ {
					exprs = append(exprs, exprStack.Pop())
				}
				// load popped exprs into the ops expr
				opsExpr.ROperand = exprs[0]
				if len(exprs) > 1 {
					opsExpr.LOperand = exprs[1]
				}
				// update the stacks
				var betterBeAnOpenParen = opsStack.Pop()
				if betterBeAnOpenParen.Op != "(" {
					return UnbalancedParenError
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
		// oh noooo!
		return DockerlangSyntaxError
	}
	c.StackTree.AST = exprStack.Pop()

	return nil
}
