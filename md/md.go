package md

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type MarkdownInfo struct {
	Slug  string
	Title string
	Date  time.Time
	Tags  []string
}

type Markdown struct {
	Info    MarkdownInfo
	Header  string
	Content []byte
}

func (m *Markdown) parseFrontmatter(r io.Reader) ([]byte, error) {
	rest, err := frontmatter.Parse(r, &m.Info)
	return rest, err
}

func (m *Markdown) parseMarkdown(b []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(b)

	htmlFlags := html.CommonFlags
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(m.extractHeader(doc), renderer)
}

func (m *Markdown) extractHeader(doc ast.Node) ast.Node {
	var heading *ast.Heading = nil
	ast.WalkFunc(doc, func(node ast.Node, entering bool) ast.WalkStatus {
		if heading != nil {
			if n := node.AsLeaf(); n != nil {
				m.Header = string(n.Literal)
				if heading != nil {
					ast.RemoveFromTree(heading)
				}
			}
			return ast.Terminate
		}
		if h, ok := node.(*ast.Heading); ok {
			if h.Level == 1 {
				heading = h
			}
		}
		return ast.GoToNext
	})
	return doc
}

func ParseMarkdownString(s string) (Markdown, error) {
	r := strings.NewReader(s)
	return parse(r)
}

func ParseMarkdownFile(filename string) (Markdown, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Markdown{}, fmt.Errorf("open file failed: %w", err)
	}
	defer f.Close()
	return parse(f)
}

func parse(r io.Reader) (Markdown, error) {
	m := Markdown{}
	rest, err := m.parseFrontmatter(r)
	if err != nil {
		return m, fmt.Errorf("parse frontmatter: %w", err)
	}

	m.Content = m.parseMarkdown(rest)

	return m, nil
}
