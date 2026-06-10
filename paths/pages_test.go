package paths

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseYaml(t *testing.T) {
	y := `paths:
 - label: Home
   file: index.md
   path: /`
	r := strings.NewReader(y)
	paths := make([]PageData, 1)
	paths[0] = PageData{Label: "Home", File: "index.md", Path: "/"}
	p, err := parsePathsYaml(r)
	assert.NoError(t, err)
	assert.Equal(t, paths, p)
}
