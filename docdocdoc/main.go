package main

import (
	"flag"
	"fmt"

	"github.com/connorwalsh/dockerlang"
)

// docdocdoc is the dockerlang official compiler.
//
// usage

const (
	UsageStr = `
docdocdoc- the official dockerlang compiler. its really the best.

usage:

    docdocdoc [flags] soucefile
`
)

func main() {
	var (
		err error
	)

	conf := &dockerlang.Config{}

	err = ParseArgs(conf)
	if err != nil {
		fmt.Println(err.Error())

		fmt.Println(UsageStr)

		return
	}

	compterpreter := dockerlang.NewCompterpreter(conf)
	err = compterpreter.Compterpret()
	if err != nil {
		fmt.Println(err.Error())

		return
	}
}

func ParseArgs(c *dockerlang.Config) error {
	// define flags
	flag.BoolVar(&c.ShowUsage, "help", false, UsageStr)

	// parse from cmd line
	flag.Parse()

	// get src filename
	args := flag.Args()

	// ensure that a source file has been provided
	if len(args) == 0 {
		return fmt.Errorf("no source file has been provided :(")
	}

	// assume the first arg is the source file
	c.SrcFileName = args[0]

	// validate flags
	if c.ShowUsage {
		fmt.Println(UsageStr)

		return nil
	}

	return nil
}
