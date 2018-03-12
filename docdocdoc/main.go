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
	// parse commandline flags
	conf := &dockerlang.Config{}
	ParseArgs(conf)

	// beging lexing and parsing
	err := dockerlang.Compile(conf)
	if err != nil {
		panic(err)
	}
}

func ParseArgs(c *dockerlang.Config) error {
	var (
		err error
	)

	// define flags
	flag.BoolVar(&c.ShowUsage, "help", false, UsageStr)

	// parse from cmd line
	flag.Parse()

	// get src filename
	args := flag.Args()

	// ensure that a source file has been provided
	if len(args) == 0 {
		err = fmt.Errorf("no source file has been provided :(")

		fmt.Println(err.Error())
		fmt.Println(UsageStr)

		return err
	}

	// assume the first arg is the source file
	c.SrcFileName = args[0]

	// validate flags
	if c.ShowUsage {
		fmt.Println(UsageStr)
	}

	return nil
}
