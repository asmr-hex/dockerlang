package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/connorwalsh/dockerlang"
)

type Computation struct {
	Type         string
	Value        string
	Dependencies []string
}

func NewComputation() *Computation {
	c := &Computation{
		Type:         os.Getenv(dockerlang.COMPUTATION_TYPE_ENV_VAR),
		Value:        os.Getenv(dockerlang.COMPUTATION_VALUE_ENV_VAR),
		Dependencies: []string{},
	}

	fmt.Printf(
		"Dockerlang Memory Ready for Use (%s %s)...\n",
		c.Type,
		c.Value,
	)

	// get the computation type and value
	depsEnvVar := os.Getenv(dockerlang.COMPUTATION_DEPS_ENV_VAR)
	computationDependencyDlcis := strings.Split(depsEnvVar, ",")
	// strings.Split returns an array with an empty string in it if there
	// are no dependcies, which is baaaad
	if computationDependencyDlcis[0] == "" {
		computationDependencyDlcis = []string{}
	}

	// get all the values this computation depends on
	for _, dlci := range computationDependencyDlcis {
		fmt.Println("dep   ", dlci)
		c.Dependencies = append(
			c.Dependencies,
			c.ExecuteDependency(dlci),
		)
	}

	return c
}

func main() {
	// make a new cozy computation
	c := NewComputation()
	http.HandleFunc("/", c.ExecHandler)
	http.HandleFunc("/kill", c.KillHandler)

	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%s", dockerlang.MEMORY_PORT),
		nil,
	))
}

func (c *Computation) ExecHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(c.Type, "exec handler")
	switch c.Type {
	case dockerlang.INT: // TODO refactor this into a constant
		io.WriteString(w, c.Value)
		fmt.Println("expect an int", c.Type, c.Value)

		//os.Exit(0)
	case dockerlang.ADDITION_OPERATOR:
		input1, _ := strconv.Atoi(c.Dependencies[0])
		input2, _ := strconv.Atoi(c.Dependencies[1])

		result := strconv.Itoa(input1 + input2)
		fmt.Println("expect addition", c.Type, result)
		io.WriteString(w, result)

		//os.Exit(0)
	default:
		// wtf is this type??????
		panic("heeellllpppp")
	}
}

func (c *Computation) KillHandler(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}

func (c *Computation) ExecuteDependency(dlci string) string {
	var (
		resultString string
	)

	fmt.Println("okay pretty basic  ", dlci)
	resp, err := http.Get(
		fmt.Sprintf("http://%s:%s/", dlci, dockerlang.MEMORY_PORT),
	)
	fmt.Println("shoulda got the resp")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		resultString = string(bodyBytes)
	}

	return resultString
}
