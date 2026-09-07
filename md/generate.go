package md

import (
	"fmt"
	"io"
	"os"

	"go.yaml.in/yaml/v4"
)

func WriteMarkdownFile(path string, info MarkdownInfo) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return writeMarkdown(f, info)
}

func writeMarkdown(w io.Writer, info MarkdownInfo) error {
	_, err := fmt.Fprintln(w, "---")
	if err != nil {
		return err
	}
	y := yaml.NewEncoder(w)
	y.SetIndent(2)
	err = y.Encode(info)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "---\n\n# %s\n", info.Title)
	if err != nil {
		return err
	}
	return nil
}
