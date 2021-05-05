package main

import (
	"bytes"
	"fmt"

	"encoding/json"

	"os/exec"

	"github.com/dsoprea/go-exif-knife"
	"github.com/dsoprea/go-logging/v2"
)

var (
	assetsPath  = ""
	appFilepath = ""
)

func CommandGetExif(filepath string, extra ...string) (exifInfo map[string]map[string]interface{}) {
	parts := []string{
		"go", "run", appFilepath, "read",
		"--filepath", filepath,
		"--json",
	}

	parts = append(parts, extra...)

	output, err := RunCommand(parts...)
	log.PanicIf(err)

	exifInfo = make(map[string]map[string]interface{})

	err = json.Unmarshal(output, &exifInfo)
	log.PanicIf(err)

	return exifInfo
}

func CommandGetExifText(filepath string, extra ...string) (exifRaw string) {
	parts := []string{
		"go", "run", appFilepath, "read",
		"--filepath", filepath,
	}

	parts = append(parts, extra...)

	output, err := RunCommand(parts...)
	log.PanicIf(err)

	return string(output)
}

func RunCommand(command_parts ...string) (output []byte, err error) {
	cmd := exec.Command(command_parts[0], command_parts[1:]...)

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
	assetsPath = exifknife.GetTestAssetsPath()
	appFilepath = exifknife.GetAppFilepath()
}
