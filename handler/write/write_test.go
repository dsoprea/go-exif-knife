package exifknifewrite

import (
	"os"
	"path"
	"testing"

	"io/ioutil"

	"github.com/dsoprea/go-jpeg-image-structure"
	"github.com/dsoprea/go-logging"
)

func TestExifWrite_Write_Noop(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintError(err)
		}
	}()

	imageFilepath := path.Join(assetsPath, "gps.jpg")

	// Write out without changes.

	ew := new(ExifWrite)

	f, err := ioutil.TempFile("", "go-exif-knife--write_test")
	log.PanicIf(err)

	outputFilepath := f.Name()

	defer os.Remove(outputFilepath)

	setTagPhrases := make([]string, 0)
	err = ew.Write(imageFilepath, setTagPhrases, outputFilepath)
	log.PanicIf(err)

	// Parse.

	jmp := jpegstructure.NewJpegMediaParser()

	sl, err := jmp.ParseFile(outputFilepath)
	log.PanicIf(err)

	rootIfd, _, err := sl.Exif()
	log.PanicIf(err)

	// Verify initial value.

	results, err := rootIfd.FindTagWithName("Software")
	log.PanicIf(err)

	if len(results) != 1 {
		t.Fatalf("'Software' tag not correctly found (1): %v", results)
	}

	value, err := rootIfd.TagValue(results[0])
	log.PanicIf(err)

	valueString := value.(string)

	if valueString != "GIMP 2.8.20" {
		t.Fatalf("Initial 'Software' tag value not correct: (%d) [%v]", len(valueString), valueString)
	}

	// Write with an update.

	setTagPhrases = []string{
		"IFD,Software,abc",
	}

	err = ew.Write(imageFilepath, setTagPhrases, outputFilepath)
	log.PanicIf(err)

	// Parse.

	sl, err = jmp.ParseFile(outputFilepath)
	log.PanicIf(err)

	rootIfd, _, err = sl.Exif()
	log.PanicIf(err)

	// Verify initial value.

	results, err = rootIfd.FindTagWithName("Software")
	log.PanicIf(err)

	if len(results) != 1 {
		t.Fatalf("'Software' tag not correctly found (2): %v", results)
	}

	value, err = rootIfd.TagValue(results[0])
	log.PanicIf(err)

	valueString = value.(string)

	if valueString != "abc" {
		t.Fatalf("Updated 'Software' tag value not correct: (%d) [%v]", len(valueString), valueString)
	}
}
