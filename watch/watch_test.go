package watch

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var tempFiles = make(map[string]bool)

func createDir(t *testing.T, path string, name string) string {
	if len(path) == 0 {
		path = "/tmp/watchit"
	}

	finalName := strings.Join([]string{path, name}, "/")
	if err := os.MkdirAll(finalName, 0700); err != nil {
		t.Fatalf(
			"failed to create temp dir; path=%s name=%s err=%s",
			path,
			name,
			err,
		)
	}

	tempFiles[finalName] = true
	return finalName
}

func createFile(t *testing.T, dir string, name string) *os.File {
	finalName := strings.Join([]string{dir, name}, "/")

	file, err := os.Create(finalName)
	if err != nil {
		t.Fatalf(
			"failed to create temp file; dir=%s name=%s err=%s",
			dir,
			name,
			err,
		)
	}

	tempFiles[finalName] = true
	return file
}

func removeFile(t *testing.T, name string) {
	if err := os.Remove(name); err != nil {
		t.Fatalf("failed to remove file: %s err=%s", name, err)
	}

	delete(tempFiles, name)
}

func writeToFile(t *testing.T, file *os.File, data string) {
	_, err := file.WriteString(data)
	if err != nil {
		t.Fatalf(
			"failed to write to temp file: file=%s data=%s err=%s",
			file.Name(),
			data,
			err,
		)
	}
}

func renameFile(t *testing.T, oldName, newName string) {
	if err := os.Rename(oldName, newName); err != nil {
		t.Fatalf(
			"failed to rename file: oldname=%s newname=%s err=%s",
			oldName,
			newName,
			err,
		)
	}

	delete(tempFiles, oldName)
	tempFiles[newName] = true
}

func cleanup(t *testing.T) {
	for file := range tempFiles {
		if err := os.RemoveAll(file); err != nil {
			t.Fatalf("failed to remove temp file: %s err=%s", file, err)
		}
	}

	tempFiles = make(map[string]bool)
}

func TestFileOps(t *testing.T) {
	defer cleanup(t)

	dir := createDir(t, "", "test_basic")
	file := createFile(t, dir, "test.md")
	defer file.Close()

	lines := []string{"hello, tests!\n", "bye, tests!\n"}

	for _, line := range lines {
		writeToFile(t, file, line)
		time.Sleep(10 * time.Millisecond)
	}

	buf := bufio.NewReader(file)
	for _, line := range lines {
		readLine, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				continue
			}
			t.Fatalf(
				"reading data from file failed; name=%s err=%s",
				file.Name(),
				err,
			)
		}
		assert.Equal(t, line, readLine)
	}
}
