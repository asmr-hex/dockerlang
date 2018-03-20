package dockerlang

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
