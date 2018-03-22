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
	Eval() (DLCI, error)
	Execute() (DLCI, error)
	GetChildren() []AST
}

// we should never actually instantiate this on its own
// it should *always* be an embedded structure (this is a mixin)
type BaseAST struct {
	ExecData *ExecutionData
}

func (b *BaseAST) GetChildren() []AST {
	return []AST{}
}

func (b *BaseAST) Execute() (DLCI, error) {
	// actual docker

	// but wait, what parts of the AST actually have their own containers?
	// containers are like memory cells and DLCIs are like memory pointers,
	// so anything we are storing in memory we will give its own docker container
	// this means variables, expressions, and functions (which are basically just
	// expressions) will have their own containers. Literals will not. The variable
	// declaration operator will spin up an new unbound variable container while
	// the variable assignment operator will set the value within that container.

	// we will be using the docker golang api since docker is written
	// in go and therefore we can just talk directly to docker through
	// this code.

	// OK before execution of anything, create a docker network

	// pass computation data (some data structure we haven't decided on yet)
	// for it to execute
	return executer.Run(b.ExecData)
}

func (b *BaseAST) Eval() (DLCI, error) {
	var (
		err error
	)

	// we need to evaluate all child expressions in order
	// to evaluate the current expression. So, evaluate
	// all the child ASTs from left to right
	for _, child := range b.GetChildren() {
		dlci, err := child.Eval()
		if err != nil {
			return "", err
		}

		// construct operands for execution
		b.ExecData.Operands = append(b.ExecData.Operands, dlci)
	}

	// we've computed all dependencies, now lets eval this thang
	literal, err := b.Execute()
	if err != nil {
		return "", err
	}

	return literal, nil
}

type Expr struct {
	BaseAST
	Name     string
	Op       string
	DLCI     string
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

func (l *Literal) Eval() (DLCI, error) {

	return l.Execute()
}

type DLCI string
