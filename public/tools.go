package public

import (
	"os"
	"os/exec"
	"path/filepath"
)

var ApplicationRoot = func() string {
	ApplicationDir, err := os.Getwd()
	if err != nil {
		file, _ := exec.LookPath(os.Args[0])
		ApplicationPath, _ := filepath.Abs(file)
		ApplicationDir, _ = filepath.Split(ApplicationPath)
	}

	return ApplicationDir
}()
