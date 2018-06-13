package main

import (
    "testing"
    "path"
    "os"
    "reflect"

    "io/ioutil"
    "encoding/json"

    "github.com/dsoprea/go-logging"
    "github.com/dsoprea/go-exif-knife"
)

var (
    appFilepath = ""
)

func getExif(filepath string) (exifInfo map[string]map[string]interface{}) {
    parts := []string {
        "go", "run", appFilepath, "read",
        "--filepath", filepath,
        "--json",
    }

    output, err := exifknife.RunCommand(parts...)
    log.PanicIf(err)

    exifInfo = make(map[string]map[string]interface{})

    err = json.Unmarshal(output, &exifInfo)
    log.PanicIf(err)

    return exifInfo
}

func TestMain_ReadAndWrite(t *testing.T) {
    imageFilepath := path.Join(assetsPath, "gps.jpg")


    // Check original value.

    exifInfo := getExif(imageFilepath)

    if reflect.DeepEqual(exifInfo["ifd0"]["Software"], "GIMP 2.8.20") != true {
        t.Fatalf("Updated 'Software' value not correct: %v", exifInfo["ifd0"]["Software"])
    }


    // Configure output file.

    f, err := ioutil.TempFile("", "go-exif-knife--write_test")
    log.PanicIf(err)

    outputFilepath := f.Name()

    defer os.Remove(outputFilepath)


    // Update the EXIF information.

    parts := []string {
        "go", "run", appFilepath, "write",
        "--filepath", imageFilepath,
        "--output-filepath", outputFilepath,
        "--set-tag", "ifd0,Software,abc",
    }

    output, err := exifknife.RunCommand(parts...)
    log.PanicIf(err)

    if len(output) != 0 {
        t.Fatalf("Expected no output:\n%s", string(output))
    }


    // Check updated value.

    exifInfo = getExif(outputFilepath)

    if reflect.DeepEqual(exifInfo["ifd0"]["Software"], "abc") != true {
        t.Fatalf("Updated 'Software' value not correct: %v", exifInfo["ifd0"]["Software"])
    }
}

func init() {
    goPath := os.Getenv("GOPATH")

    appFilepath = path.Join(goPath, "src", "github.com", "dsoprea", "go-exif-knife", "command", "go-exif-knife", "main.go")
}
