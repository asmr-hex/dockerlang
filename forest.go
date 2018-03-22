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

func (b *BaseAST) Execute() error {
	// actual docker

	// but wait, what parts of the AST actually have their own containers?
	// containers are like memory cells and DLIIs are like memory pointers,
	// so anything we are storing in memory we will give its own docker container
	// this means variables, expressions, and functions (which are basically just
	// expressions) will have their own containers. Literals will not. The variable
	// declaration operator will spin up an new unbound variable container while
	// the variable assignment operator will set the value within that container.
	//
	// we will be using the docker golang api since docker is written
	// in go and therefore we can just talk directly to docker through
	// this code.
	// before execution of anything, create a docker network
	// when we run a docker container, we must give it the network name
	// and the computation type (some data structure we haven't decided on yet)
	// for it to execute
	// each container will persist until that computation is no longer in scope
	// within the source code.
	// to that end, we can implement a very simple garbage collector or something.
	return nil
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
	err = b.Execute()
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
