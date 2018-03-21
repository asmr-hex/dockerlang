package dockerlang

// each function in Dockerlang will have its own StackTree of the form:
//
//              function StackTree
//        _______________|_______________
//       |          |          |         |
//      args      global      local    body
//      /\          /\         /\        |
//    AST map    AST map  empty AST map  AST

type StackTree struct {
	Name    string
	Args    map[string]AST
	Locals  map[string]AST
	Globals map[string]AST
	Body    AST
}

func NewStackTree(name string) *StackTree {
	return &StackTree{
		Name: name,
	}
}

func (s *StackTree) Eval() ([]interface{}, []interface{}, error) {
	return nil, nil, nil
}

type Expr struct {
	Name     string
	Op       string
	DLII     string
	Arity    int
	Operands []interface{}
}

func (e *Expr) Eval() ([]interface{}, []interface{}, error) {
	return nil, nil, nil
}

type AST interface {
	Eval() ([]interface{}, []interface{}, error)
}
