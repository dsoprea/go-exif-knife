package main

import (
	"os"
	"path"
	"reflect"
	"testing"

	"io/ioutil"

	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-exif-knife"
)

var (
	appFilepath = ""
)

func TestMain_Read(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	// Check original value.

	exifInfo := exifknife.CommandGetExif(imageFilepath)

	if reflect.DeepEqual(exifInfo["ifd0"]["Software"], "GIMP 2.8.20") != true {
		t.Fatalf("'Software' value not correct: %v", exifInfo["ifd0"]["Software"])
	}
}

func TestMain_Write(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	// Check original value.

	exifInfo := exifknife.CommandGetExif(imageFilepath)

	if reflect.DeepEqual(exifInfo["ifd0"]["Software"], "GIMP 2.8.20") != true {
		t.Fatalf("Updated 'Software' value not correct: %v", exifInfo["ifd0"]["Software"])
	}

	// Configure output file.

	f, err := ioutil.TempFile("", "go-exif-knife--write_test")
	log.PanicIf(err)

	outputFilepath := f.Name()

	defer os.Remove(outputFilepath)

	// Update the EXIF information.

	parts := []string{
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

	exifInfo = exifknife.CommandGetExif(outputFilepath)

	if reflect.DeepEqual(exifInfo["ifd0"]["Software"], "abc") != true {
		t.Fatalf("Updated 'Software' value not correct: %v", exifInfo["ifd0"]["Software"])
	}
}

func TestMain_Gps(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	parts := []string{
		"go", "run", appFilepath, "gps",
		"--filepath", imageFilepath,
		"--json",
	}

	output, err := exifknife.RunCommand(parts...)
	log.PanicIf(err)

	expected := `{
    "Altitude": 0,
    "LatitudeDecimal": 26.586666666666666,
    "LongitudeDecimal": -80.05361111111111,
    "Timestamp": "2018-04-29T01:22:57Z",
    "TimestampUnix": 1524964977
}
`

	if string(output) != expected {
		t.Fatalf("GPS result not correct.")
	}
}

func init() {
	goPath := os.Getenv("GOPATH")

	appFilepath = path.Join(goPath, "src", "github.com", "dsoprea", "go-exif-knife", "command", "go-exif-knife", "main.go")
}
