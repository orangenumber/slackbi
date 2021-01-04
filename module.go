package slackbi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

type Modules struct {
	p           *SBI
	ModuleDir   string // Directory where module stays
	ConfNam     string
	Items       map[string]*Module
	lastRefresh time.Time
}

func (m *Modules) Load() error {
	if m.p == nil {
		return m_sbi_ptr_error
	}
	modDir := m.p.config.Module.Dir
	files, err := ioutil.ReadDir(modDir)
	if err != nil {
		return err
	}

	// MUTEX RW LOCK.
	m.p.mu.Lock()
	defer m.p.mu.Unlock()

	// CLEAR ITEMS MAP
	m.Items = make(map[string]*Module)
	m.p.modules.ModuleDir = modDir

	var loadedModules []string // used for logging

	for _, v := range files {
		if v.IsDir() {
			// read config files
			mo, err := m.readModuleItemConf(v.Name())
			if err != nil {
				return err
			}

			// Make sure entry point file has runnable permission
			finfo, err := os.Stat(path.Join(modDir, v.Name(), mo.Module.EntryPoint))
			if err != nil {
				return err
			}

			if CheckFilePerm(finfo.Mode()) == false {
				err = fmt.Errorf("permission error, file <%s>", path.Join(modDir, v.Name(), mo.Module.EntryPoint))
				return err
			}

			// All went successful
			m.Items[v.Name()] = mo
			loadedModules = append(loadedModules, v.Name())
		}
	}

	sort.Strings(loadedModules)

	m.p.logger.Infof("Loaded %d module(s), module=%s", len(m.Items), strings.Join(loadedModules, ","))
	m.lastRefresh = time.Now()
	return nil
}

func (m *Modules) readModuleItemConf(filename string) (*Module, error) {
	filename = path.Join(m.p.config.Module.Dir, filename, m.p.configFilename)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	var c Module
	dec := json.NewDecoder(f)
	if err := dec.Decode(&c); err != nil {
		return nil, err
	}
	c.dir = path.Dir(filename)
	return &c, nil
}

func (m *Modules) Get(moduleName string) (*Module, error) {
	if v, ok := m.Items[moduleName]; ok {
		return v, nil
	} else {
		return nil, mf_mod_not_found_sName
	}
}
