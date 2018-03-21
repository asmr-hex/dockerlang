package dockerlang

// each function in Dockerlang will have its own StackTree of the form:
//
//                 function StackTree
//        __________|____________________
//       |          |          |         |
//      arg       global      local    returns
//      /\          /\         /\       /\
//     ASTs        ASTs       ASTs     ASTs

type StackTree struct {
	Name    string
	Args    []AST
	Local   []AST
	Global  []AST
	Returns []AST
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
	Op       string
	DLII     string
	Arity    int
	LOperand interface{}
	ROperand interface{}
}

func (e *Expr) Eval() ([]interface{}, []interface{}, error) {
	return nil, nil, nil
}

type AST interface {
	Eval() ([]interface{}, []interface{}, error)
}
