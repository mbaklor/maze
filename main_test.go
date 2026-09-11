package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	var s string
	c := ""
	e := ""
	f := ""
	merge("test1", &s, c, e)
	assert.Equal(t, "test1", s)

	c = ""
	e = ""
	f = ""
	merge("test1", &s, c, e, f)
	assert.Equal(t, "test1", s)

	c = ""
	e = ""
	f = "test2"
	merge("test1", &s, c, e, f)
	assert.Equal(t, "test2", s)

	c = ""
	e = "test3"
	f = "test2"
	merge("test1", &s, c, e, f)
	assert.Equal(t, "test3", s)

	c = "test4"
	e = "test3"
	f = "test2"
	merge("test1", &s, c, e, f)
	assert.Equal(t, "test4", s)

	var si int
	ci := 0
	ei := 0
	fi := 0
	merge(1, &si, ci, ei, fi)
	assert.Equal(t, 1, si)

	ci = 0
	ei = 0
	fi = 2
	merge(1, &si, ci, ei, fi)
	assert.Equal(t, 2, si)

	ci = 0
	ei = 3
	fi = 2
	merge(1, &si, ci, ei, fi)
	assert.Equal(t, 3, si)

	ci = 4
	ei = 3
	fi = 2
	merge(1, &si, ci, ei, fi)
	assert.Equal(t, 4, si)
}
