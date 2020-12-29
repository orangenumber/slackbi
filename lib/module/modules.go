package module

import (
	"fmt"
	"github.com/gonyyi/aface"
	"io/ioutil"
	"os"
	"path"
)

func NewModules(logger aface.Logger1a, dir string, confName string) (*Modules, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, v := range files {
		if v.IsDir() {
			// read config files
			mo, err := ReadModuleFile(path.Join(dir, v.Name(), confName))
			if err != nil {
				return nil, err
			}

			// Make sure entry point file has runnable permission
			finfo, err := os.Stat(path.Join(dir, v.Name(), mo.Module.EntryPoint))
			if err != nil {
				return nil, err
			}

			if CheckFilePerm(finfo.Mode()) == false {
				err = fmt.Errorf("permission error, file <%s>", path.Join(dir, v.Name(), mo.Module.EntryPoint))
				return nil, err
			}
		}
	}

	m := &Modules{ModuleDir: dir, ConfNam: confName}
	if logger != nil {
		m.logger = logger
	} else {
		m.logger = &aface.LoggerDummy1a{}
	}
	return m, nil
}

type Modules struct {
	logger    aface.Logger1a
	ModuleDir string // Directory where module stays
	ConfNam   string
}

// IsExist check for a directory
func (m *Modules) IsExist(moduleName string) bool {
	fileinfo, err := os.Stat(path.Join(m.ModuleDir, moduleName))
	if err != nil {
		return false
	}
	if fileinfo.IsDir() {
		return true
	}
	return false
}

func (m *Modules) Get(moduleName string) (*Module, error) {
	mod, err := ReadModuleFile(path.Join(m.ModuleDir, moduleName, m.ConfNam))
	if err != nil {
		return nil, err
	}

	// Check if entry point exists
	fi, err := os.Stat(path.Join(m.ModuleDir, moduleName, mod.Module.EntryPoint))
	if err != nil {
		return nil, err
	}
	if fi.IsDir() == true {
		return nil, fmt.Errorf("entry_point cannot be a directory, entry_point=%s", fi.Name())
	}

	return mod, err
}
