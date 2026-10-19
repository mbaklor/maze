package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"strings"
)

func UsageFunc(s string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Maze: dynamic markdown static page server\nversion: %s\n\n%s\n", readVersion(), s)
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
			vsplit := strings.FieldsFunc(info.Main.Version, func(r rune) bool {
				return r == '-' || r == '+'
			})
			if len(vsplit) == 1 {
				return vsplit[0]
			}
			v := vsplit[0]
			hash := ""
			dirty := ""
			for _, s := range info.Settings {
				if s.Key == "vcs.revision" {
					hash = "-" + s.Value[:7]
				}
				if s.Key == "vcs.modified" && s.Value == "true" {
					dirty = "-modified"
				}
			}
			return v + hash + dirty
		}
	}
	return version
}

func printVersion() {
	fmt.Fprintf(os.Stderr, "Maze version: %s\n", readVersion())
}
