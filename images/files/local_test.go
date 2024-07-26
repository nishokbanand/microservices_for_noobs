package files

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupLocal(t *testing.T) (*Local, string, func()) {
	dir, err := os.MkdirTemp("", "temp")
	if err != nil {
		t.Fatal("cannot create temp dir", err)
	}
	store, err := NewLocal(dir, 1024*1000*5)
	if err != nil {
		t.Fatal(err)
	}
	return store, dir, func() {
		os.RemoveAll("temp")
	}
}

func TestSavesContentOfReader(t *testing.T) {
	savePath := "1/test.png"
	fileContent := "hello world"
	l, dir, cleanup := setupLocal(t)
	defer cleanup()
	err := l.Save(savePath, bytes.NewBuffer([]byte(fileContent)))
	assert.NoError(t, err)
	f, err := os.Open(filepath.Join(dir, savePath))
	assert.NoError(t, err)
	d, err := io.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, fileContent, string(d))
}
