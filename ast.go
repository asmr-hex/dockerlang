package dockerlang

// each function in Dockerlang will have its own StackTree of the form:
//
//                 function
//               ______|______
//               |           |
//            scope         body
//            ___|__         |
//           |      |      return
//          args  vars       /\
//           /\    /\       Nodes
//         Nodes  Nodes

type StackTree struct {
	Name  string
	Scope *Scope
	AST   *AST
}

func NewStackTree(name string) *StackTree {
	return &StackTree{
		Name:  name,
		Scope: NewScope(),
		AST:   NewAST(),
	}
}

type Scope struct {
	Args []Node
	Vars []Node
}

func NewScope() *Scope {
	return &Scope{
		Args: []Node{},
		Vars: []Node{},
	}
}

type AST struct {
	Return Node
}

func NewAST() *AST {
	return &AST{}
}

type Node struct {
	DLII     string // (DLII = dockerlang image identifier)
	Children []Node
}
