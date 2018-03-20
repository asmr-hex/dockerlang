package dockerlang

// each function in Dockerlang will have its own StackTree of the form:
//
//                 function StackTree
//                   ______|______
//                  |             |
//                Scope          AST
//            ______|______
//           |      |      |
//          args  global  local
//           /\     /\     /\
//          ASTs   ASTs   ASTs

type StackTree struct {
	Name  string
	Scope *Scope
	AST   Evaluable
}

func NewStackTree(name string) *StackTree {
	return &StackTree{
		Name:  name,
		Scope: NewScope(),
	}
}

func (s *StackTree) Eval() ([]interface{}, []interface{}, error) {
	return nil, nil, nil
}

type Scope struct {
	Args []Evaluable
	Vars []Evaluable
}

func NewScope() *Scope {
	return &Scope{
		Args: []Evaluable{},
		Vars: []Evaluable{},
	}
}

type Expr struct {
	Op       string
	DLII     string
	Arity    int
	lOperand interface{}
	rOperand interface{}
}

func (e *Expr) Eval() ([]interface{}, []interface{}, error) {
	return nil, nil, nil
}

type Evaluable interface {
	Eval() ([]interface{}, []interface{}, error)
}
