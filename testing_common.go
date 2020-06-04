package exifknife

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"encoding/json"
	"go/build"
	"os/exec"

	"github.com/dsoprea/go-logging"
)

var (
	assetsPath  = ""
	appFilepath = ""
)

// GetModuleRootPath returns our source-path when running from source during
// tests.
func GetModuleRootPath() string {
	p, err := build.Default.Import(
		"github.com/dsoprea/go-exif-knife",
		build.Default.GOPATH,
		build.FindOnly)

	log.PanicIf(err)

	packagePath := p.Dir
	return packagePath
}

// CommandGetExif calls the reader subcommand.
func CommandGetExif(filepath string) (exifInfo map[string]map[string]interface{}) {
	parts := []string{
		"go", "run", appFilepath, "read",
		"--filepath", filepath,
		"--json",
	}

	output, err := RunCommand(parts...)
	log.PanicIf(err)

	exifInfo = make(map[string]map[string]interface{})

	err = json.Unmarshal(output, &exifInfo)
	log.PanicIf(err)

	return exifInfo
}

// RunCommand calls an arbitrary command and captures the output.
func RunCommand(commandParts ...string) (output []byte, err error) {
	cmd := exec.Command(commandParts[0], commandParts[1:]...)

	b := new(bytes.Buffer)
	cmd.Stdout = b
	cmd.Stderr = b

	err = cmd.Run()
	raw := b.Bytes()

	if err != nil {
		fmt.Printf(string(raw))
		log.Panic(err)
	}

	return b.Bytes(), nil
}

// getModuleRootPath returns our source-path when running from source during
// tests.
func getModuleRootPath() string {
	moduleRootPath := os.Getenv("EXIF_KNIFE_MODULE_ROOT_PATH")
	if moduleRootPath != "" {
		return moduleRootPath
	}

	currentWd, err := os.Getwd()
	log.PanicIf(err)

	currentPath := currentWd
	visited := make([]string, 0)

	for {
		tryStampFilepath := path.Join(currentPath, ".MODULE_ROOT")

		_, err := os.Stat(tryStampFilepath)
		if err != nil && os.IsNotExist(err) != true {
			log.Panic(err)
		} else if err == nil {
			break
		}

		visited = append(visited, tryStampFilepath)

		currentPath = path.Dir(currentPath)
		if currentPath == "/" {
			log.Panicf("could not find module-root: %v", visited)
		}
	}

	return currentPath
}

// GetTestAssetsPath returns the test-asset path.
func GetTestAssetsPath() string {
	if assetsPath == "" {
		moduleRootPath := getModuleRootPath()
		assetsPath = path.Join(moduleRootPath, "assets")
	}

	return assetsPath
}

// GetAppFilepath returns the file-path of the main command.
func GetAppFilepath() string {
	moduleRootPath := getModuleRootPath()
	appFilepath = path.Join(moduleRootPath, "command", "go-exif-knife", "main.go")

	return appFilepath
}
