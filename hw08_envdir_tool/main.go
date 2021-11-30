package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		log.Fatal("There aren't enough arguments")
	}

	envDir, cmd := args[1], args[2:]

	environment, err := ReadDir(envDir)
	if err != nil {
		log.Fatal(err)
	}
	RunCmd(cmd, environment)
}
