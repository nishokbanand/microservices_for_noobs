package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	basePath string
	maxLimit int
}

func NewLocal(basePath string, maxLimit int) (*Local, error) {
	fp, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	fmt.Println("abs path", fp, "normal_path", basePath)
	return &Local{basePath, maxLimit}, nil
}

func (l *Local) Save(path string, contents io.Reader) error {
	lp := l.fullPath(path)
	dir := filepath.Dir(lp)
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("Could not create directory %v", err)
	}
	_, err = os.Stat(lp)
	if err == nil {
		err = os.Remove(lp)
		if err != nil {
			return fmt.Errorf("Could not delete file %v", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("Could not get the file %v", err)
	}
	f, err := os.Create(lp)
	if err != nil {
		return fmt.Errorf("Could not create the file %v", err)
	}
	defer f.Close()
	_, err = io.Copy(f, contents)
	if err != nil {
		return fmt.Errorf("Could not copy contents to the file %v", err)
	}
	return nil
}

func (l *Local) fullPath(path string) string {
	return filepath.Join(l.basePath, path)
}
