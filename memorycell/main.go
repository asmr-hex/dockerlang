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

var (
	COMPUTATION_TYPE             string
	COMPUTATION_VALUE            string
	COMPUTATION_DEPENDENCY_DLCIS []string
	COMPUTATION_DEPENDENCIES     []string
)

func main() {

	// get the computation type and value
	COMPUTATION_TYPE = os.Getenv(dockerlang.COMPUTATION_TYPE_ENV_VAR)
	COMPUTATION_VALUE = os.Getenv(dockerlang.COMPUTATION_VALUE_ENV_VAR)
	depsEnvVar := os.Getenv(dockerlang.COMPUTATION_DEPS_ENV_VAR)
	COMPUTATION_DEPENDENCY_DLCIS = strings.Split(depsEnvVar, ",")

	// get all the values this computation depends on
	COMPUTATION_DEPENDENCIES = []string{}
	for _, dlci := range COMPUTATION_DEPENDENCY_DLCIS {
		if dlci != "" {
			COMPUTATION_DEPENDENCIES = append(COMPUTATION_DEPENDENCIES, ExecuteDependency(dlci))
		}
	}

	http.HandleFunc("/", ExecHandler)
	http.HandleFunc("/kill", KillHandler)

	fmt.Printf(
		"Dockerlang Memory Ready for Use (%s %s)...\n",
		COMPUTATION_TYPE,
		COMPUTATION_VALUE,
	)
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%s", dockerlang.MEMORY_PORT),
		nil,
	))
}

func ExecHandler(w http.ResponseWriter, r *http.Request) {
	switch COMPUTATION_TYPE {
	case dockerlang.INT: // TODO refactor this into a constant
		io.WriteString(w, COMPUTATION_VALUE)

		os.Exit(0)
	case dockerlang.ADDITION_OPERATOR:
		input1, _ := strconv.Atoi(COMPUTATION_DEPENDENCIES[0])
		input2, _ := strconv.Atoi(COMPUTATION_DEPENDENCIES[1])

		result := input1 + input2
		io.WriteString(w, string(result))

		os.Exit(0)
	default:
		// wtf is this type??????
		panic("heeellllpppp")
	}

	fmt.Println("EXEC SOMETHING")
}

func KillHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("KILL ME")
	os.Exit(0)
}

func ExecuteDependency(dlci string) string {
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
