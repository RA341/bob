package bob

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RA341/bob/parser"
	"github.com/RA341/bob/vm"
)

// FileName todo add lowercase support
const FileName = "Bobfile"

func Run(cmd string, args []string) error {
	return fmt.Errorf("bobfile parsing is unimplemented")

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	bobFilePath := getBobFilePath(wd)
	//bobWorkingDir := filepath.Dir(bobFilePath)

	ins, err := parser.ParseBobFromPath(bobFilePath)
	if err != nil {
		log.Fatalf("Failed to parse Bobfile: %v", err)
	}

	vmm := new(vm.VM)
	vmm.Start(ins, vm.DefaultFns)

	return nil
}

func AddDefaultEnvs(workingDir string) {
	//b.Vars.Add("OS", runtime.GOOS)
	//b.Vars.Add("ARCH", runtime.GOARCH)
	//b.Vars.Add("WorkDir", workingDir)
}

func getBobFilePath(bobBase string) string {
	bobFilePath := ""

	var pathsTried []string
	for {
		if bobBase == "/" {
			break
		}

		tmpPath := filepath.Join(bobBase, FileName)
		stat, err := os.Stat(tmpPath)
		if err == nil {
			if stat.IsDir() {
				log.Fatalf("%s is a directory", bobFilePath)
			}

			bobFilePath = tmpPath
			break
		}

		if !os.IsNotExist(err) {
			log.Fatalln("Failed to stat bobfile:", err)
		}
		pathsTried = append(pathsTried, bobBase)
		bobBase = filepath.Dir(bobBase)
	}

	if bobFilePath == "" {
		join := strings.Join(pathsTried, "\n")
		log.Fatalf("Failed to find Bobfile, tried the following paths \n%s", join)
	}

	return bobFilePath
}
