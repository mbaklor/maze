package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
)

func UsageFunc(s string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Maze: dynamic markdown static page server\nversion: %s\n\n%s\n", version, s)
	}
}

var version = "develop"

const usage = `USAGE
maze [option] [flags]

OPTIONS
serve			start the HTTP server
generate		create a new page file with the expected frontmatter fields
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

func readVersion() string {
	if version != "develop" {
		return version
	}
	info, ok := debug.ReadBuildInfo()
	if ok {
		if info.Main.Version != "(devel)" && info.Main.Version != "" {
			return info.Main.Version
		}
	}
	return version
}

func printVersion() {
	fmt.Fprintf(os.Stderr, "Maze version: %s\n", readVersion())
}
