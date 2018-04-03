package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/connorwalsh/dockerlang"
)

type Computation struct {
	Type         string
	Value        string
	Dependencies []string
	Stop         chan struct{}
}

func NewComputation() *Computation {
	c := &Computation{
		Type:         os.Getenv(dockerlang.COMPUTATION_TYPE_ENV_VAR),
		Value:        os.Getenv(dockerlang.COMPUTATION_VALUE_ENV_VAR),
		Dependencies: []string{},
		Stop:         make(chan struct{}),
	}

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
	http.HandleFunc("/", c.HttpExecHandler)
	http.HandleFunc("/kill", c.KillHandler)

	// exit is a special case where we automatically run the http response
	if c.Type == dockerlang.EXIT_OPERATOR {
		fmt.Println(c.ExecHandler())

		return
	}

	go c.ListenForTermination()

	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%s", dockerlang.MEMORY_PORT),
		nil,
	))
}

// we need this because the http handlers need to return and we can't
// directly os.Exit within them.
func (c *Computation) ListenForTermination() {
	select {
	case <-c.Stop:
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}
}

func (c *Computation) ExecHandler() string {
	switch c.Type {
	case dockerlang.EXIT_OPERATOR:
		return c.Dependencies[0]
	case dockerlang.INT:
		return c.Value
	case dockerlang.ADDITION_OPERATOR:
		input1, _ := strconv.Atoi(c.Dependencies[0])
		input2, _ := strconv.Atoi(c.Dependencies[1])

		result := strconv.Itoa(input1 + input2)
		return result
	case dockerlang.SUBTRACTION_OPERATOR:
		input1, _ := strconv.Atoi(c.Dependencies[0])
		input2, _ := strconv.Atoi(c.Dependencies[1])

		result := strconv.Itoa(input1 - input2)
		return result
	case dockerlang.MULTIPLICATION_OPERATOR:
		input1, _ := strconv.Atoi(c.Dependencies[0])
		input2, _ := strconv.Atoi(c.Dependencies[1])

		result := strconv.Itoa(input1 * input2)
		return result
	case dockerlang.DIVISION_OPERATOR:
		input1, _ := strconv.Atoi(c.Dependencies[0])
		input2, _ := strconv.Atoi(c.Dependencies[1])

		result := strconv.Itoa(input1 / input2)
		return result
	case dockerlang.MODULO_OPERATOR:
		input1, _ := strconv.ParseFloat(c.Dependencies[0], 32)
		input2, _ := strconv.ParseFloat(c.Dependencies[1], 32)

		result := strconv.FormatFloat(
			math.Mod(input1, input2),
			'f',
			-1,
			32,
		)

		return result
	default:
		// wtf is this type??????
		panic("heeellllpppp")

	}
}

func (c *Computation) HttpExecHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, c.ExecHandler())

	switch c.Type {
	// case dockerlang.VARIABLE:
	default:
		c.Stop <- struct{}{}
	}

}

func (c *Computation) KillHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	os.Exit(0)
}

func (c *Computation) ExecuteDependency(dlci string) string {
	var (
		resultString string
	)

	resp, err := http.Get(
		fmt.Sprintf("http://%s:%s/", dlci, dockerlang.MEMORY_PORT),
	)
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
