package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "mispipe: Wrong number of args, aborting")
		os.Exit(1)
	}
	cmds := make([]*exec.Cmd, len(os.Args)-1)
	for i, arg := range os.Args[1:] {
		args := []string{"/bin/sh", "-c"}
		if runtime.GOOS == "windows" {
			args = []string{"cmd", "/c"}
		}
		args = append(args, arg)
		cmds[i] = exec.Command(args[0], args[1:]...)
		if i > 0 {
			rc, err := cmds[i-1].StdoutPipe()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s: %s\n", os.Args[0], arg, err)
				os.Exit(1)
			}
			cmds[i].Stdin = rc
		}
	}
	cmds[0].Stdin = os.Stdin
	cmds[len(cmds)-1].Stdout = os.Stdout
	for _, c := range cmds {
		err := c.Start()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
			os.Exit(1)
		}
	}
	for _, c := range cmds {
		err := c.Wait()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		}
	}
	status := 0
	state := cmds[0].ProcessState
	if state != nil {
		if ws, ok := state.Sys().(syscall.WaitStatus); ok {
			if ws.Exited() {
				status = ws.ExitStatus()
			}
		}
	}
	os.Exit(status)
}
