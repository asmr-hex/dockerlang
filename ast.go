package dockerlang

// each function in Dockerlang will have its own AST of the form:
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

type AST struct {
	Name  string
	Scope *Scope
	Body  *Body
}

type Scope struct {
	Args []Node
	Vars []Node
}

type Body struct {
	Return Node
}

type Node struct {
	DLII     string // (DLII = dockerlang image identifier)
	Children []Node
}
