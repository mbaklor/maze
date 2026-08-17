package md_test

import (
	"os"
	"testing"
	"time"

	"github.com/mbaklor/website/md"
	"github.com/stretchr/testify/assert"
)

const markdown = `---
title: test title
slug: test-slug
date: 2026-03-25
tags:
  - one
  - two
  - three
---

# test heading

test body

## test subheading

- one
- two
`

const html = `<p>test body</p>

<h2 id="test-subheading">test subheading</h2>

<ul>
<li>one</li>
<li>two</li>
</ul>
`

func TestParseMarkdownFile(t *testing.T) {
	_, err := md.ParseMarkdownFile("nope.md")
	assert.ErrorIs(t, err, os.ErrNotExist)
	m, err := md.ParseMarkdownFile("test.md")
	assert.NoError(t, err)
	assert.Equal(t, "test title", m.Info.Title)
	assert.Equal(t, "test-slug", m.Info.Slug)
	assert.Equal(t, time.Date(2026, time.March, 25, 0, 0, 0, 0, time.UTC), m.Info.Date)
	assert.Contains(t, m.Info.Tags, "one")
	assert.Contains(t, m.Info.Tags, "two")
	assert.Contains(t, m.Info.Tags, "three")

	assert.Equal(t, "test heading", m.Header)
	assert.Equal(t, html, string(m.Content))
}

func TestParseMarkdownString(t *testing.T) {
	m, err := md.ParseMarkdownString(markdown)
	assert.NoError(t, err)
	assert.Equal(t, "test title", m.Info.Title)
	assert.Equal(t, "test-slug", m.Info.Slug)
	assert.Equal(t, time.Date(2026, time.March, 25, 0, 0, 0, 0, time.UTC), m.Info.Date)
	assert.Contains(t, m.Info.Tags, "one")
	assert.Contains(t, m.Info.Tags, "two")
	assert.Contains(t, m.Info.Tags, "three")

	assert.Equal(t, "test heading", m.Header)
	assert.Equal(t, html, string(m.Content))
}
