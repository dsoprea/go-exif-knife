package exifknife

import (
	"path"
	"reflect"
	"testing"

	"github.com/dsoprea/go-exif"
	"github.com/dsoprea/go-logging"
)

func TestGetExif_Jpeg(t *testing.T) {
	filepath := path.Join(assetsPath, "image.jpg")

	mc, err := GetExif(filepath)
	log.PanicIf(err)

	if mc.MediaType != "jpeg" {
		t.Fatalf("media-type not 'jpeg'")
	}

	ti := exif.NewTagIndex()

	it, err := ti.GetWithName(exif.IfdPathStandard, "Model")
	log.PanicIf(err)

	ite := mc.RootIfd.EntriesByTagId[it.Id][0]

	value, err := mc.RootIfd.TagValue(ite)
	log.PanicIf(err)

	expected := "Canon EOS 5D Mark III"
	if value.(string) != expected {
		t.Fatalf("model not valid")
	}
}

func TestGetExif_Png(t *testing.T) {
	filepath := path.Join(assetsPath, "image.png")

	mc, err := GetExif(filepath)
	log.PanicIf(err)

	if mc.MediaType != "png" {
		t.Fatalf("media-type not 'png'")
	}

	ti := exif.NewTagIndex()

	it, err := ti.GetWithName(exif.IfdPathStandard, "ImageWidth")
	log.PanicIf(err)

	ite := mc.RootIfd.EntriesByTagId[it.Id][0]

	value, err := mc.RootIfd.TagValue(ite)
	log.PanicIf(err)

	expected := []uint32{11}
	if reflect.DeepEqual(value, expected) != true {
		t.Fatalf("image-width not valid")
	}
}

func TestGetExif_Other(t *testing.T) {
	filepath := path.Join(assetsPath, "image.tiff")

	mc, err := GetExif(filepath)
	log.PanicIf(err)

	if mc.MediaType != "other" {
		t.Fatalf("media-type not 'other'")
	}

	ti := exif.NewTagIndex()

	it, err := ti.GetWithName(exif.IfdPathStandard, "Artist")
	log.PanicIf(err)

	ite := mc.RootIfd.EntriesByTagId[it.Id][0]

	value, err := mc.RootIfd.TagValue(ite)
	log.PanicIf(err)

	expected := "Jean Cornillon"
	if value.(string) != expected {
		t.Fatalf("artist not correct")
	}
}
