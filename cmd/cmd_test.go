package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	dirs := []string{"/tmp"}
	cmds := []string{"ls"}
	ext := []string{"md"}
	expected := new(dirs, cmds, ext, false, false)

	os.Args = []string{
		"-d", "/tmp",
		"-c", "ls",
	}
	actual := Parse()

	assert.Equal(
		t,
		expected,
		actual,
		"expected and actual results should be same",
	)
}

func TestMultipleDirs(t *testing.T) {
	dirs := []string{"/tmp", "~/Downloads"}
	cmds := []string{"ls -lh"}
	ext := []string{"txt"}
	expected := new(dirs, cmds, ext, false, false)

	os.Args = []string{
		"-d", "/tmp",
		"-d", "~/Downloads",
		"-c", "ls -lh",
		"-e", "txt",
	}
	actual := Parse()

	assert.Equal(
		t,
		expected,
		actual,
		"expected and actual results should be same",
	)
}

func TestMultipleCmds(t *testing.T) {
	dirs := []string{"~/Downloads"}
	cmds := []string{"ls -lh", "finger"}
	ext := []string{"html"}
	expected := new(dirs, cmds, ext, false, false)

	os.Args = []string{
		"-d", "~/Downloads",
		"-c", "ls -lh",
		"-c", "finger",
		"-e", "html",
	}
	actual := Parse()

	assert.Equal(
		t,
		expected,
		actual,
		"expected and actual results should be same",
	)
}

func TestMultipleExts(t *testing.T) {
	dirs := []string{"~/Downloads"}
	cmds := []string{"ls -lh", "finger"}
	ext := []string{"md", "txt"}
	expected := new(dirs, cmds, ext, false, false)

	os.Args = []string{
		"-d", "~/Downloads",
		"-c", "ls -lh",
		"-c", "finger",
		"-e", "md",
		"-e", "txt",
	}
	actual := Parse()

	assert.Equal(
		t,
		expected,
		actual,
		"expected and actual results should be same",
	)
}

func TestRegex(t *testing.T) {
	const (
		regexMd    = ".+\\.(md)$"
		regexTxt   = ".+\\.(txt)$"
		regexHTML  = ".+\\.(html)$"
		regexMdTxt = ".+\\.(md|txt)$"
	)

	var regexPairs = []struct {
		ext   []string
		regex string
	}{
		{[]string{"md"}, regexMd},
		{[]string{"txt"}, regexTxt},
		{[]string{"html"}, regexHTML},
		{[]string{"md", "txt"}, regexMdTxt},
	}

	for _, r := range regexPairs {
		assert.Equal(t, r.regex, buildRegex(r.ext))
	}
}
