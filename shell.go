package main

import "os/exec"

type IShellCommand interface {
	Run() error
	Output() ([]byte, error)
}

type execCmd struct {
	*exec.Cmd
}

func shellCommandWrapper(name string, arg ...string) IShellCommand {
	return execCmd{Cmd: exec.Command(name, arg...)}
}

var shellCommand = shellCommandWrapper
