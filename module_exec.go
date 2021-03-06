package slackbi

import (
	"bytes"
	"os/exec"
	"path"
)

// This is old legacy version 1
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

// interface_version 2
func (m *Module) ExecV2(jsonData []byte) (jsonOutput []byte, errorReceived []byte, err error) {
	var c *exec.Cmd

	c = exec.Command("./" + m.Module.EntryPoint)
	// pwd, _ := os.Getwd()
	c.Dir = path.Join(m.dir)

	c.Stdin = bytes.NewBuffer(jsonData)

	var sout bytes.Buffer
	var serr bytes.Buffer
	c.Stdout = &sout
	c.Stderr = &serr

	err = c.Run()

	return sout.Bytes(), serr.Bytes(), err
}
