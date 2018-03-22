package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/connorwalsh/dockerlang"
)

var (
	COMPUTATION_TYPE  string
	COMPUTATION_VALUE string
)

func main() {

	// get the computation type and value
	COMPUTATION_TYPE = os.Getenv(dockerlang.COMPUTATION_TYPE_ENV_VAR)
	COMPUTATION_VALUE = os.Getenv(dockerlang.COMPUTATION_VALUE_ENV_VAR)

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
	fmt.Println("EXEC SOMETHING")
}

func KillHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("KILL ME")
}
