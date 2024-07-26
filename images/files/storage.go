package files

import "io"

type Storage interface {
	Save(path string, contents io.Reader) error
}
