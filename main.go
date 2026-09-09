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

FLAGS
-c	--config	Path to server config file (default: "config.yml")

GLOBAL FLAGS
-h	--help		Print this usage message`

func main() {
	flag.Usage = UsageFunc(usage)

	//need to make sure I have the flag in case we run with no command
	var configPath string
	flag.StringVar(&configPath, "config", "config.yml", "path to server config file")
	flag.StringVar(&configPath, "c", "config.yml", "path to server config file")

	flag.Parse()
	if len(flag.Args()) == 0 {
		serve(os.Args[1:])
		return
	}
	command := os.Args[1]
	switch command {
	case "serve":
		serve(os.Args[2:])
	case "version":
		printVersion()
	default:
		flag.Usage()
	}
}

func printVersion() {
	fmt.Fprintf(os.Stderr, "maze version: %s\n", version)
}
