package exifknife

import (
	"io"
	"os"
	"path"
	"reflect"
	"testing"

	"io/ioutil"

	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"github.com/dsoprea/go-logging/v2"
)

func TestGetExif_Jpeg(t *testing.T) {
	assetsPath := GetTestAssetsPath()
	filepath := path.Join(assetsPath, "image.jpg")

	mc, err := GetExif(filepath)
	log.PanicIf(err)

	if mc.MediaType != JpegMediaType {
		t.Fatalf("Media-type not correct for a JPEG.")
	}

	ti := exif.NewTagIndex()

	it, err := ti.GetWithName(exifcommon.IfdStandardIfdIdentity, "Model")
	log.PanicIf(err)

	ite := mc.RootIfd.EntriesByTagId()[it.Id][0]

	value, err := ite.Value()
	log.PanicIf(err)

	expected := "Canon EOS 5D Mark III"
	if value.(string) != expected {
		t.Fatalf("Model not valid.")
	}
}

func TestGetExif_Png(t *testing.T) {
	assetsPath := GetTestAssetsPath()
	filepath := path.Join(assetsPath, "image.png")

	mc, err := GetExif(filepath)
	log.PanicIf(err)

	if mc.MediaType != PngMediaType {
		t.Fatalf("Media-type not correct for a PNG.")
	}

	ti := exif.NewTagIndex()

	it, err := ti.GetWithName(exifcommon.IfdStandardIfdIdentity, "ImageWidth")
	log.PanicIf(err)

	ite := mc.RootIfd.EntriesByTagId()[it.Id][0]

	value, err := ite.Value()
	log.PanicIf(err)

	expected := []uint32{11}
	if reflect.DeepEqual(value, expected) != true {
		t.Fatalf("Image-width not valid.")
	}
}

func TestGetExif_Heic(t *testing.T) {
	assetsPath := GetTestAssetsPath()
	filepath := path.Join(assetsPath, "image.heic")

	mc, err := GetExif(filepath)
	log.PanicIf(err)

	if mc.MediaType != HeicMediaType {
		t.Fatalf("Media-type not correct for an HEIC.")
	}

	ti := exif.NewTagIndex()

	it, err := ti.GetWithName(exifcommon.IfdStandardIfdIdentity, "XResolution")
	log.PanicIf(err)

	ite := mc.RootIfd.EntriesByTagId()[it.Id][0]

	value, err := ite.Value()
	log.PanicIf(err)

	expected := []exifcommon.Rational{{Numerator: 72, Denominator: 1}}
	if reflect.DeepEqual(value, expected) != true {
		t.Fatalf("Image-width not valid: %v", value)
	}
}

func TestGetExif_Tiff(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintError(err)

			t.Fatalf("Test failure.")
		}
	}()

	assetsPath := GetTestAssetsPath()
	filepath := path.Join(assetsPath, "image.tiff")

	mc, err := GetExif(filepath)
	log.PanicIf(err)

	if mc.MediaType != TiffMediaType {
		t.Fatalf("Media-type not correct for an TIFF.")
	}

	ti := exif.NewTagIndex()

	it, err := ti.GetWithName(exifcommon.IfdStandardIfdIdentity, "Artist")
	log.PanicIf(err)

	ite := mc.RootIfd.EntriesByTagId()[it.Id][0]

	value, err := ite.Value()
	log.PanicIf(err)

	expected := "Jean Cornillon"
	if value.(string) != expected {
		t.Fatalf("Artist not correct.")
	}
}

func TestGetExif_Other(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintError(err)

			t.Fatalf("Test failure.")
		}
	}()

	f, err := ioutil.TempFile("", "")
	log.PanicIf(err)

	copyFilepath := f.Name()

	defer func() {
		f.Close()

		os.Remove(copyFilepath)
	}()

	assetsPath := GetTestAssetsPath()
	originalFilepath := path.Join(assetsPath, "image.tiff")

	g, err := os.Open(originalFilepath)
	log.PanicIf(err)

	// Write some padding at the top of the file so that it's no longer a
	// recognized file-type.

	blank := make([]byte, 4)

	_, err = f.Write(blank)
	log.PanicIf(err)

	// Now, copy the image verbatim.

	_, err = io.Copy(f, g)
	log.PanicIf(err)

	// Try the reason.

	mc, err := GetExif(copyFilepath)
	log.PanicIf(err)

	if mc.MediaType != "other" {
		t.Fatalf("Media-type not 'other' as expected.")
	}

	ti := exif.NewTagIndex()

	it, err := ti.GetWithName(exifcommon.IfdStandardIfdIdentity, "Artist")
	log.PanicIf(err)

	ite := mc.RootIfd.EntriesByTagId()[it.Id][0]

	value, err := ite.Value()
	log.PanicIf(err)

	expected := "Jean Cornillon"
	if value.(string) != expected {
		t.Fatalf("Artist not correct.")
	}
}
