package md_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mbaklor/maze/md"
	"github.com/stretchr/testify/assert"
)

const genFrontmatter = `---
slug: test-slug
title: test title
date: "2026-03-25"
tags:
  - one
  - two
  - three
---

# test title
`

func TestGenerateMarkdown(t *testing.T) {
	m, err := md.ParseMarkdownString(markdown)
	assert.NoError(t, err)

	tempfile := filepath.Join(t.TempDir(), "out.md")
	err = md.WriteMarkdownFile(tempfile, m.Info)
	assert.NoError(t, err)

	f, err := os.ReadFile(tempfile)
	assert.NoError(t, err)

	assert.Equal(t, genFrontmatter, string(f))
}
