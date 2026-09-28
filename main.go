package main

import (
	"flag"
	"fmt"
	"os"
)

func UsageFunc(s string) func() {
	return func() {
		fmt.Fprintln(os.Stderr, s)
	}
}

const version = "0.0.2"

const usage = `Maze: dynamic markdown static page server
version ` + version + `

USAGE
maze [option] [flags]

OPTIONS
serve			start the HTTP server
version			print application version

GLOBAL FLAGS
-h	--help		Print this usage message`

func main() {
	flag.Usage = UsageFunc(usage)

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}
	command := os.Args[1]
	switch command {
	case "serve":
		serve(os.Args[2:])
	case "generate":
		generate(os.Args[2:])
	case "version":
		printVersion()
	default:
		flag.Usage()
	}
}

func printVersion() {
	fmt.Fprintf(os.Stderr, "maze version: %s\n", version)
}
