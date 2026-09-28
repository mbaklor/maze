package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mbaklor/maze/md"
)

type generateConfig struct {
	md.MarkdownInfo
}

const generateUsage = `Maze: dynamic markdown static page server
version ` + version + `

USAGE
maze generate [flags] <page-path>

REQUIRED ARGUMENT
filename		Name of the file to create, must have either .md extension or none

FLAGS
-t	--title		Post/page title (required)
-d	--description	Post/page description
-p	--page-title	HTML title for post/page (default: same as post title)
-g	--tag		Any amount of tags to add to the post/page

GLOBAL FLAGS
-h	--help		Print this usage message`

func parseGenerateFlags(args []string) generateConfig {
	var gc generateConfig
	set := flag.NewFlagSet("generate", flag.ExitOnError)
	set.Usage = UsageFunc(generateUsage)

	set.StringVar(&gc.PageTitle, "page-title", "", "")
	set.StringVar(&gc.PageTitle, "p", "", "")

	set.StringVar(&gc.Description, "description", "", "")
	set.StringVar(&gc.Description, "d", "", "")

	set.StringVar(&gc.Title, "title", "", "")
	set.StringVar(&gc.Title, "t", "", "")

	set.Func("tag", "", func(s string) error {
		gc.Tags = append(gc.Tags, s)
		return nil
	})
	set.Func("g", "", func(s string) error {
		gc.Tags = append(gc.Tags, s)
		return nil
	})

	set.Parse(args)
	gc.Slug = set.Arg(0)

	if gc.Slug == "" {
		fmt.Fprint(os.Stderr, "need a path for file\n\n")
		set.Usage()
		os.Exit(1)
	}

	if gc.Title == "" {
		fmt.Fprint(os.Stderr, "title is required!\n\n")
		set.Usage()
		os.Exit(1)
	}

	gc.Date = md.CurrentDate()
	return gc
}

func interactiveGenerateConfig() generateConfig {
	var gc generateConfig
	fmt.Println("Not yet implmented")
	return gc
}

func generate(args []string) {
	var gc generateConfig
	if len(args) == 0 {
		gc = interactiveGenerateConfig()
	} else {
		gc = parseGenerateFlags(args)
	}
	fmt.Printf("date: %v, now: %v\n", gc.Date, time.Now().Local())
}
