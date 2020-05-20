package exifknife

import (
	"path"
	"reflect"
	"testing"

	"github.com/dsoprea/go-exif/v2"
	"github.com/dsoprea/go-exif/v2/common"
	"github.com/dsoprea/go-logging"
)

func TestGetExif_Jpeg(t *testing.T) {
	filepath := path.Join(assetsPath, "image.jpg")

	mc, err := GetExif(filepath)
	log.PanicIf(err)

	if mc.MediaType != JpegMediaType {
		t.Fatalf("Media-type not correct for a JPEG.")
	}

	ti := exif.NewTagIndex()

	it, err := ti.GetWithName(exifcommon.IfdStandardIfdIdentity.UnindexedString(), "Model")
	log.PanicIf(err)

	ite := mc.RootIfd.EntriesByTagId[it.Id][0]

	value, err := ite.Value()
	log.PanicIf(err)

	expected := "Canon EOS 5D Mark III"
	if value.(string) != expected {
		t.Fatalf("Model not valid.")
	}
}

func TestGetExif_Png(t *testing.T) {
	filepath := path.Join(assetsPath, "image.png")

	mc, err := GetExif(filepath)
	log.PanicIf(err)

	if mc.MediaType != PngMediaType {
		t.Fatalf("Media-type not correct for a PNG.")
	}

	ti := exif.NewTagIndex()

	it, err := ti.GetWithName(exifcommon.IfdStandardIfdIdentity.UnindexedString(), "ImageWidth")
	log.PanicIf(err)

	ite := mc.RootIfd.EntriesByTagId[it.Id][0]

	value, err := ite.Value()
	log.PanicIf(err)

	expected := []uint32{11}
	if reflect.DeepEqual(value, expected) != true {
		t.Fatalf("Image-width not valid.")
	}
}

func TestGetExif_Heic(t *testing.T) {
	filepath := path.Join(assetsPath, "image.heic")

	mc, err := GetExif(filepath)
	log.PanicIf(err)

	if mc.MediaType != HeicMediaType {
		t.Fatalf("Media-type not correct for an HEIC.")
	}

	ti := exif.NewTagIndex()

	it, err := ti.GetWithName(exifcommon.IfdStandardIfdIdentity.UnindexedString(), "XResolution")
	log.PanicIf(err)

	ite := mc.RootIfd.EntriesByTagId[it.Id][0]

	value, err := ite.Value()
	log.PanicIf(err)

	expected := []exifcommon.Rational{{Numerator: 72, Denominator: 1}}
	if reflect.DeepEqual(value, expected) != true {
		t.Fatalf("Image-width not valid: %v", value)
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

	filepath := path.Join(assetsPath, "image.tiff")

	mc, err := GetExif(filepath)
	log.PanicIf(err)

	if mc.MediaType != "other" {
		t.Fatalf("Media-type not 'other' as expected.")
	}

	ti := exif.NewTagIndex()

	it, err := ti.GetWithName(exifcommon.IfdStandardIfdIdentity.UnindexedString(), "Artist")
	log.PanicIf(err)

	ite := mc.RootIfd.EntriesByTagId[it.Id][0]

	value, err := ite.Value()
	log.PanicIf(err)

	expected := "Jean Cornillon"
	if value.(string) != expected {
		t.Fatalf("Artist not correct.")
	}
}
