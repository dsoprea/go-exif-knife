package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"io/ioutil"

	"github.com/dsoprea/go-logging"
)

func TestMain_Read_JustTry(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	// Check original value.

	mediaType := CommandGetExifText(imageFilepath, "--just-try")
	mediaType = strings.TrimSpace(mediaType)

	if mediaType != "jpeg" {
		t.Fatalf("'just-try' didn't work for JPEG: [%s]", mediaType)
	}
}

func TestMain_Read(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	// Check original value.

	exifInfo := CommandGetExif(imageFilepath)

	if reflect.DeepEqual(exifInfo["IFD"]["Software"], "GIMP 2.8.20") != true {
		t.Fatalf("'Software' value not correct: %v", exifInfo["IFD"]["Software"])
	}
}

func TestMain_Read_Text(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	// Check original value.

	exifRaw := CommandGetExifText(imageFilepath, "--ifd", "IFD1")

	expected :=
		` IFD: Ifd<ID=(3) IFD-PATH=[IFD] INDEX=(1) COUNT=(6) OFF=(0x039e) CHILDREN=(0) PARENT=(0x0000) NEXT-IFD=(0x0000)>
 - TAG: IfdTagEntry<TAG-IFD-PATH=[IFD1] TAG-ID=(0x0103) TAG-TYPE=[SHORT] UNIT-COUNT=(1)> NAME=[Compression] VALUE=[[6]]
 - TAG: IfdTagEntry<TAG-IFD-PATH=[IFD1] TAG-ID=(0x011a) TAG-TYPE=[RATIONAL] UNIT-COUNT=(1)> NAME=[XResolution] VALUE=[[72/1]]
 - TAG: IfdTagEntry<TAG-IFD-PATH=[IFD1] TAG-ID=(0x011b) TAG-TYPE=[RATIONAL] UNIT-COUNT=(1)> NAME=[YResolution] VALUE=[[72/1]]
 - TAG: IfdTagEntry<TAG-IFD-PATH=[IFD1] TAG-ID=(0x0128) TAG-TYPE=[SHORT] UNIT-COUNT=(1)> NAME=[ResolutionUnit] VALUE=[[2]]
`

	if exifRaw != expected {
		fmt.Printf("ACTUAL:\n")
		fmt.Printf("\n")
		fmt.Println(exifRaw)
		fmt.Printf("\n")

		fmt.Printf("EXPECTED:\n")
		fmt.Printf("\n")
		fmt.Println(expected)
		fmt.Printf("\n")

		t.Fatalf("IFD-specific read not correct")
	}
}

func TestMain_Read_SpecificIfd(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	// Check original value.

	exifRaw := CommandGetExifText(imageFilepath, "--ifd", "IFD1")

	expected :=
		` IFD: Ifd<ID=(3) IFD-PATH=[IFD] INDEX=(1) COUNT=(6) OFF=(0x039e) CHILDREN=(0) PARENT=(0x0000) NEXT-IFD=(0x0000)>
 - TAG: IfdTagEntry<TAG-IFD-PATH=[IFD1] TAG-ID=(0x0103) TAG-TYPE=[SHORT] UNIT-COUNT=(1)> NAME=[Compression] VALUE=[[6]]
 - TAG: IfdTagEntry<TAG-IFD-PATH=[IFD1] TAG-ID=(0x011a) TAG-TYPE=[RATIONAL] UNIT-COUNT=(1)> NAME=[XResolution] VALUE=[[72/1]]
 - TAG: IfdTagEntry<TAG-IFD-PATH=[IFD1] TAG-ID=(0x011b) TAG-TYPE=[RATIONAL] UNIT-COUNT=(1)> NAME=[YResolution] VALUE=[[72/1]]
 - TAG: IfdTagEntry<TAG-IFD-PATH=[IFD1] TAG-ID=(0x0128) TAG-TYPE=[SHORT] UNIT-COUNT=(1)> NAME=[ResolutionUnit] VALUE=[[2]]
`

	if exifRaw != expected {
		fmt.Printf("ACTUAL:\n")
		fmt.Printf("\n")
		fmt.Println(exifRaw)
		fmt.Printf("\n")

		fmt.Printf("EXPECTED:\n")
		fmt.Printf("\n")
		fmt.Println(expected)
		fmt.Printf("\n")

		t.Fatalf("IFD-specific read not correct")
	}
}

func TestMain_Read_SpecificTag(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	// Check original value.

	exifRaw := CommandGetExifText(imageFilepath, "--tag", "ResolutionUnit", "--json")

	expected :=
		`{
    "IFD": {
        "ResolutionUnit": [
            2
        ]
    },
    "IFD1": {
        "ResolutionUnit": [
            2
        ]
    }
}
`

	if strings.TrimSpace(exifRaw) != strings.TrimSpace(expected) {
		fmt.Printf("ACTUAL:\n%s\n", exifRaw)
		fmt.Printf("EXPECTED:\n%s\n", expected)

		t.Fatalf("Tag-specific read not correct")
	}
}

func TestMain_Read_SpecificIfdAndSpecificTag(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	// Check original value.

	exifRaw := CommandGetExifText(imageFilepath, "--ifd", "IFD1", "--tag", "ResolutionUnit", "--json")

	expected :=
		`{
    "IFD1": {
        "ResolutionUnit": [
            2
        ]
    }
}
`

	if strings.TrimSpace(exifRaw) != strings.TrimSpace(expected) {
		fmt.Printf("ACTUAL:\n%s\n", exifRaw)
		fmt.Printf("EXPECTED:\n%s\n", expected)

		t.Fatalf("Tag-specific read not correct")
	}
}

func TestMain_Read_JustValues(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	// Check original value.

	value := CommandGetExifText(imageFilepath, "--ifd", "IFD1", "--tag", "ResolutionUnit", "--just-values")

	expected :=
		`2
`

	if strings.TrimSpace(value) != strings.TrimSpace(expected) {
		fmt.Printf("ACTUAL:\n%s\n", value)
		fmt.Printf("EXPECTED:\n%s\n", expected)

		t.Fatalf("Just-value response not correct")
	}
}

func TestMain_Write(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	// Check original value.

	exifInfo := CommandGetExif(imageFilepath)

	if reflect.DeepEqual(exifInfo["IFD"]["Software"], "GIMP 2.8.20") != true {
		t.Fatalf("Updated 'Software' value not correct: %v", exifInfo["IFD"]["Software"])
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
		"--set-tag", "IFD,Software,abc",
	}

	output, err := RunCommand(parts...)
	log.PanicIf(err)

	if len(output) != 0 {
		t.Fatalf("Expected no output:\n%s", string(output))
	}

	// Check updated value.

	exifInfo = CommandGetExif(outputFilepath)

	if reflect.DeepEqual(exifInfo["IFD"]["Software"], "abc") != true {
		t.Fatalf("Updated 'Software' value not correct: %v", exifInfo["IFD"]["Software"])
	}
}

func TestMain_Gps(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	parts := []string{
		"go", "run", appFilepath, "gps",
		"--filepath", imageFilepath,
		"--json",
	}

	output, err := RunCommand(parts...)
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

func TestMain_Thumbnail_Write(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")
	expectedThumbnailFilepath := path.Join(assetsPath, "gps.jpg.thumbnail")

	f, err := ioutil.TempFile("", "go-exif-knife-TestMain_Thumbnail_Write")
	log.PanicIf(err)

	outputFilepath := f.Name()

	defer f.Close()
	defer os.Remove(outputFilepath)

	parts := []string{
		"go", "run", appFilepath, "thumbnail",
		"--filepath", imageFilepath,
		"--output-filepath", outputFilepath,
	}

	output, err := RunCommand(parts...)
	log.PanicIf(err)

	if string(output) != "" {
		t.Fatalf("Output not expected:\n%s\n", string(output))
	}

	actual, err := ioutil.ReadFile(outputFilepath)
	log.PanicIf(err)

	expected, err := ioutil.ReadFile(expectedThumbnailFilepath)
	log.PanicIf(err)

	if bytes.Compare(actual, expected) != 0 {
		t.Fatalf("Thumbnail not correct.")
	}
}
