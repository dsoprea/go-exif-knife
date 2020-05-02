package exifkniferead

import (
	"path"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestExifRead_Read(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	er := new(ExifRead)

	justTry := false
	specificIfdDesignation := ""
	specificTags := []string{}
	justPrintValues := false
	printAsJson := false
	ignoreNoExif := false

	err := er.Read(imageFilepath, justTry, specificIfdDesignation, specificTags, justPrintValues, printAsJson, ignoreNoExif)
	log.PanicIf(err)
}

func TestExifRead_Read_JustTry(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	er := new(ExifRead)

	justTry := true
	specificIfdDesignation := ""
	specificTags := []string{}
	justPrintValues := false
	printAsJson := false
	ignoreNoExif := false

	err := er.Read(imageFilepath, justTry, specificIfdDesignation, specificTags, justPrintValues, printAsJson, ignoreNoExif)
	log.PanicIf(err)
}

func TestExifRead_Read_SpecificIfd(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	er := new(ExifRead)

	justTry := false
	specificIfdDesignation := "IFD0"
	specificTags := []string{}
	justPrintValues := false
	printAsJson := false
	ignoreNoExif := false

	err := er.Read(imageFilepath, justTry, specificIfdDesignation, specificTags, justPrintValues, printAsJson, ignoreNoExif)
	log.PanicIf(err)
}

func TestExifRead_Read_SpecificTags(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	er := new(ExifRead)

	justTry := false
	specificIfdDesignation := ""
	specificTags := []string{"ResolutionUnit"}
	justPrintValues := false
	printAsJson := false
	ignoreNoExif := false

	err := er.Read(imageFilepath, justTry, specificIfdDesignation, specificTags, justPrintValues, printAsJson, ignoreNoExif)
	log.PanicIf(err)
}

func TestExifRead_Read_SpecificTags_InvalidIsOkay(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	er := new(ExifRead)

	justTry := false
	specificIfdDesignation := ""
	specificTags := []string{"xyz"}
	justPrintValues := false
	printAsJson := false
	ignoreNoExif := false

	err := er.Read(imageFilepath, justTry, specificIfdDesignation, specificTags, justPrintValues, printAsJson, ignoreNoExif)
	log.PanicIf(err)
}

func TestExifRead_Read_SpecificIfdAndSpecificTags(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	er := new(ExifRead)

	justTry := false
	specificIfdDesignation := "IFD"
	specificTags := []string{"ResolutionUnit"}
	justPrintValues := false
	printAsJson := false
	ignoreNoExif := false

	err := er.Read(imageFilepath, justTry, specificIfdDesignation, specificTags, justPrintValues, printAsJson, ignoreNoExif)
	log.PanicIf(err)
}

func TestExifRead_Read_Json(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	er := new(ExifRead)

	justTry := false
	specificIfdDesignation := ""
	specificTags := []string{}
	justPrintValues := false
	printAsJson := true
	ignoreNoExif := false

	err := er.Read(imageFilepath, justTry, specificIfdDesignation, specificTags, justPrintValues, printAsJson, ignoreNoExif)
	log.PanicIf(err)
}

func TestExifRead_Read_JustPrint(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	er := new(ExifRead)

	justTry := false
	specificIfdDesignation := ""
	specificTags := []string{}
	justPrintValues := true
	printAsJson := false
	ignoreNoExif := false

	err := er.Read(imageFilepath, justTry, specificIfdDesignation, specificTags, justPrintValues, printAsJson, ignoreNoExif)
	log.PanicIf(err)
}
