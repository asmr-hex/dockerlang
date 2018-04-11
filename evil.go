package dockerlang

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types/container"
)

// receive a stack tree that only has a body AST set
// Evaluation:
// traverse the Body AST like a beautiful red sailboat, maintain reference to parent StackTree for lookups in scope
// switch node:
// ≡ -> assign name to locals, create docker image
// = -> lookup in Locals/Args/Globals and reassign value to Docker image
// ...

// Precedence for variable lookup is:
// Local, Args, Global

func (c *Compterpreter) Evaluate() error {
	b, _ := json.MarshalIndent(c.StackTree, "    ", "  ")
	fmt.Println(string(b))

	/*
		for _, operand := range c.StackTree.Operands {
			b, _ := json.MarshalIndent(operand, "    ", "  ")
			fmt.Println(string(b))
		}
	*/
	r, err := c.StackTree.Operands[0].Eval()

	wait, errChan := executer.Docker.ContainerWait(
		context.Background(),
		r,
		container.WaitConditionNotRunning,
	)

	select {
	case <-wait:
		// execution complete!
	case err = <-errChan:
		// uo oh
	}

	return err
}
