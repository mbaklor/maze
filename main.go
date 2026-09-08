package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "0.0.1"

const usage = `Maze: dynamic markdown static page server
version ` + version + `

USAGE
maze [option] [flags]

OPTIONS
serve			start the HTTP server
version			print application version

GLOBAL FLAGS
--help	-h		Print this usage message`

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
	}
	if len(os.Args) == 1 {
		serve()
		return
	}
	command := os.Args[1]
	switch command {
	case "serve":
		serve()
	case "version":
		printVersion()
	default:
		flag.Usage()
	}
}

func printVersion() {
	fmt.Fprintf(os.Stderr, "maze version: %s\n", version)
}
