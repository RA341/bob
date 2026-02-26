package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// BobFileName todo maybe add lowercase support
const BobFileName = "Bobfile"

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
	}

	bobFilePath := getBobFilePath(wd)
	//log.Printf("Bobfile path: %s", bobFilePath)

	var bobFile Bobfile
	ParseFromFile(&bobFile, bobFilePath)

	Execute(&bobFile)
}

func getBobFilePath(bobBase string) string {
	bobFilePath := ""

	var pathsTried []string
	for {
		if bobBase == "/" {
			break
		}

		tmpPath := filepath.Join(bobBase, BobFileName)
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
