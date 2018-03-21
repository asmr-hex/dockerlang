package dockerlang

// each function in Dockerlang will have its own StackTree of the form:
//
//              function StackTree
//        _______________|_______________
//       |          |          |         |
//      args      global      local   body AST
//      /\          /\         /\        \
//    AST map    AST map  empty AST map   (AST of parsed code)

type AST interface {
	Eval() error
	GetChildren() []AST
}

// we should never actually instantiate this on its own
// it should *always* be an embedded structure (this is a mixin)
type BaseAST struct{}

func (b *BaseAST) GetChildren() []AST {
	return []AST{}
}

func (b *BaseAST) Eval() error {
	var (
		err error
	)

	// we need to evaluate all child expressions in order
	// to evaluate the current expression. So, evaluate
	// all the child ASTs from left to right
	for _, child := range b.GetChildren() {
		err = child.Eval()
		if err != nil {
			return err
		}
	}

	// we've computed all dependencies, now lets eval this thang
	err = b.Eval()
	if err != nil {
		return err
	}

	return nil
}

type Expr struct {
	BaseAST
	Name     string
	Op       string
	DLII     string
	Arity    int
	Operands []AST
	Args     map[string]AST
	Locals   map[string]AST
	Globals  map[string]AST
}

func NewExpr(name string) *Expr {
	return &Expr{
		Name: name,
	}
}

type IfConditional struct{}

func (c *IfConditional) Eval() error {

	return nil
}

type Variable struct {
	BaseAST
	Literal
	Name  string
	Bound bool
}

type Literal struct {
	BaseAST
	Type  string
	Value interface{}
}
