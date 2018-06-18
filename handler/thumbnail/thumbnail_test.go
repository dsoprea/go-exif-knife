package exifknifethumbnail

import (
	"bytes"
	"os"
	"path"
	"testing"

	"io/ioutil"

	"github.com/dsoprea/go-logging"
)

func TestExifThumbnail_ExtractThumbnail(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")
	expectedThumbnailFilepath := path.Join(assetsPath, "gps.jpg.thumbnail")

	f, err := ioutil.TempFile("", "go-exif-knife--TestExifThumbnail_ExtractThumbnail")
	log.PanicIf(err)

	outputFilepath := f.Name()

	defer os.Remove(outputFilepath)

	et := new(ExifThumbnail)

	err = et.ExtractThumbnail(imageFilepath, outputFilepath)
	log.PanicIf(err)

	actual, err := ioutil.ReadFile(outputFilepath)
	log.PanicIf(err)

	expected, err := ioutil.ReadFile(expectedThumbnailFilepath)
	log.PanicIf(err)

	if bytes.Compare(actual, expected) != 0 {
		t.Fatalf("Thumbnail not correct.")
	}
}

func TestExifThumbnail_writeBytes(t *testing.T) {
	f, err := ioutil.TempFile("", "go-exif-knife--TestExifThumbnail_writeBytes")
	log.PanicIf(err)

	outputFilepath := f.Name()

	defer os.Remove(outputFilepath)

	et := new(ExifThumbnail)

	data := []byte{'a', 'b', 'c'}

	err = et.writeBytes(outputFilepath, data)
	log.PanicIf(err)

	actual, err := ioutil.ReadFile(outputFilepath)
	log.PanicIf(err)

	if bytes.Compare(actual, data) != 0 {
		t.Fatalf("Bytes not written correctly.")
	}
}
