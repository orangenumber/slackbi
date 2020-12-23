package module

import (
	"encoding/json"
	"io/ioutil"
)

type Module struct {
	Module struct {
		InterfaceVersion int    `json:"interface_version"`
		EntryPoint       string `json:"entry_point"`
		AvgRuntimeSec    int    `json:"avg_runtime_sec"`
	} `json:"module"`
	Info struct {
		Name         string   `json:"name"`
		Version      string   `json:"version"`
		Created      string   `json:"created"`
		LastModified string   `json:"last-modified"`
		Contacts     []string `json:"contacts"`
	} `json:"info"`
	Help struct {
		Intro   string   `json:"intro"`
		Website string   `json:"website"`
		Help    []string `json:"help"`
		Usage   []string `json:"usage"`
	} `json:"help"`
	dir string // this is to be filled when load
}

func (m *Module) Save(filename string) error {
	b, err := json.MarshalIndent(m, "", "   ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, b, 0755)
}
