package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for k, v := range env {
		if v.NeedRemove {
			if err := os.Unsetenv(k); err != nil {
				log.Fatal(err)
			}
		}
		if err := os.Setenv(k, v.Value); err != nil {
			log.Fatal(err)
		}
	}

	command, params := cmd[0], cmd[1:]

	e := exec.Command(command, params...)
	e.Stdin = os.Stdin
	e.Stdout = os.Stdout

	if err := e.Run(); err != nil {
		log.Println(err)
	}

	return e.ProcessState.ExitCode()
}
