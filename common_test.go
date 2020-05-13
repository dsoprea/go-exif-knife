package exifknife

import (
	"bytes"
	"fmt"
	"os"
	"path"

	"encoding/json"

	"os/exec"

	"github.com/dsoprea/go-logging"
)

var (
	assetsPath  = ""
	appFilepath = ""
)

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

func init() {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		log.Panicf("GOPATH is empty")
	}

	assetsPath = path.Join(goPath, "src", "github.com", "dsoprea", "go-exif-knife", "assets")
	appFilepath = path.Join(goPath, "src", "github.com", "dsoprea", "go-exif-knife", "command", "go-exif-knife", "main.go")
}
