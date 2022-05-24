package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func WorkingDir() string {
	cwd, _ := os.Getwd()
	return cwd
}

func ExecutablePath() string {
	path, _ := os.Executable()
	return path
}

func ExecutableDir() string {
	return filepath.Dir(ExecutablePath())
}

func ResetWorkingDir() {
	cwd := WorkingDir()
	exd := ExecutableDir()

	if strings.Compare(cwd, exd) != 0 {
		os.Chdir(ExecutableDir())
	}
}
