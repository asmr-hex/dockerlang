package dockerlang

// each function in Dockerlang will have its own StackTree of the form:
//
//                          Expr
//        ___________________|___________________
//       |          |            |              |
//      args      global        local        Operands
//      /\          /\           /\             /\
//    AST map    AST map   empty AST map   (ASTs of parsed code)

// All nodes within the Abstract Syntax Tree we are parsing should be
// sub ASTs which satisfy the interface{} AST which only requires that
// we can evaluate that subtree. Importantly, the evaluation of all ASTs
// will result a DockerLang Container Id which points to a memorycell
// container which holds the value of the computation at that node.
type AST interface {
	Eval() (string, error)
}

// An Expr is an expression which satisfies the AST interface.
type Expr struct {
	Op       string
	Arity    int
	Operands []AST
	Args     map[string]AST
	Locals   map[string]AST
	Globals  map[string]AST
}

// def f():
//   1 + 1
//   return 2 + 2

// evaluate an expression
func (e *Expr) Eval() (string, error) {
	execData := &ExecutionData{
		ComputationType: e.Op,
	}

	// we need to evaluate all child expressions in order
	// to evaluate the current expression. So, evaluate
	// all the child ASTs from left to right
	for _, child := range e.Operands {
		dlci, err := child.Eval()
		if err != nil {
			return "", err
		}

		// construct operands for execution
		execData.Operands = append(execData.Operands, dlci)
	}

	// we've computed all dependencies, now lets eval this thang

	return executer.Run(execData)
}

// TODO implement me. this should embed an Expr (since it has the same fields)
// but is should overwrite the Eval function since it does that differently.
type IfConditional struct{}

/*
SOME NOTES (mr,cw|4.25.2018):

------ define function ------
(# add (x y) (+ x y))

build an expression:
{
    Op: ?
    Operands: [<sequential list of expressions in function>]
    Locals: [local vars or functions]
    Args: [x, y]
}
store that expression in Locals

------ eval program --------
(add 4 5)

build an expression for our executed program:
add
/ \
4  5


{ Op: "add"
  operands: [4, 5]
}


------

locals
add:
 +
/ \
x  y -> args

                  /|\
                  \./
              nnnnnnnnnnn
             mm0m0m0m0m0m0m
             m0m0m0m0m0m0mm
             mmmmmmmmmmmmmm
           ///            \\\
         ////              \\\\
        ////                \\\\
         |||   ====   ====  |||
          ||   (  )   (  )  ||
           |        |       |
           |        |       |
            \      ( )     /
             \   mmmmmmm  /
              |   -----  |
               \________/

               welcome 2 dockerland

*/

// an Identifier can be a named expression which can either be a variable
// or a function which is indicated by its Type.
type Identifier struct {
	Type  string
	Name  string
	Bound bool
}

func (i *Identifier) Eval() (string, error) {
	return executer.Run(
		&ExecutionData{
			ComputationType: i.Type,
		},
	)
}

type Literal struct {
	Type  string
	Value string
}

func (l *Literal) Eval() (string, error) {
	return executer.Run(
		&ExecutionData{
			ComputationType: l.Type,
			Value:           l.Value,
		},
	)
}
