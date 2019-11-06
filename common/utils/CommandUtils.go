package utils

import (
	"bytes"
	"github.com/kris-nova/logger"
	"os"
	"os/exec"
	"strings"
)

const ShellToUse = "bash"

//todo - if in debug mode write to stdout + file logger

func ShellNoExit(command string) (string ,bool)  {
	return execute(command ,"" ,false)
}

func Shell(command string) string {
	result ,_ := execute(command ,"" ,true)
	return result
}

func ShellCurrentDir(command string ,currentDir string) string {
	result ,_ := execute(command ,currentDir ,true)
	return result
}


func execute(command string ,currentDir string ,exitOnFailure bool) (string ,bool) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	logger.Always("cmd =" + command)
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if currentDir != "" {
		cmd.Dir = currentDir
	}

	err := cmd.Run()

	if err != nil {
		logger.Critical("Unable to execute command= " + command, err)
		logger.Critical("result=" + stderr.String())
		if exitOnFailure {
			os.Exit(1)
		} else {
			return stdout.String() ,false
		}
	}
	return strings.TrimSuffix(stdout.String(), "\n") ,true
}

