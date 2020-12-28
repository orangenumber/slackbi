package module

import (
	"encoding/json"
	"io"
	"os"
	"path"
)

func ReadModuleFile(filename string) (*Module, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return ReadModule(f, path.Dir(filename))
}

func ReadModule(ior io.Reader, dir string) (*Module, error) {
	var c Module
	dec := json.NewDecoder(ior)
	if err := dec.Decode(&c); err != nil {
		return nil, err
	}
	c.dir = dir
	return &c, nil
}
