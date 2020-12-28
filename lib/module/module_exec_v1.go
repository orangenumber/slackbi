package module

import (
	"bytes"
	"os/exec"
	"path"
)

// This is old legacy version
func (m *Module) ExecV1(user string, command string) (output []byte, outputError []byte, err error) {
	var c *exec.Cmd

	c = exec.Command(path.Join(m.dir, m.Module.EntryPoint), user, command)

	var sout bytes.Buffer
	var serr bytes.Buffer
	c.Stdout = &sout
	c.Stderr = &serr

	err = c.Run()

	return sout.Bytes(), serr.Bytes(), err
}
