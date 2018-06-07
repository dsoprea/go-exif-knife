package exifknife

import (
    "testing"
    "path"
    "reflect"

    "github.com/dsoprea/go-logging"
    "github.com/dsoprea/go-exif"
)

func TestGetExif_Jpeg(t *testing.T) {
    filepath := path.Join(assetsPath, "image.jpg")

    mediaType, rootIfd, err := GetExif(filepath)
    log.PanicIf(err)

    if mediaType != "jpeg" {
        t.Fatalf("media-type not 'jpeg'")
    }


    ti := exif.NewTagIndex()

    it, err := ti.GetWithName(exif.RootIi, "Model")
    log.PanicIf(err)


    ite := rootIfd.EntriesByTagId[it.Id][0]

    value, err := rootIfd.TagValue(ite)
    log.PanicIf(err)

    expected := "Canon EOS 5D Mark III"
    if value.(string) != expected {
        t.Fatalf("model not valid")
    }
}

func TestGetExif_Png(t *testing.T) {
    filepath := path.Join(assetsPath, "image.png")

    mediaType, rootIfd, err := GetExif(filepath)
    log.PanicIf(err)

    if mediaType != "png" {
        t.Fatalf("media-type not 'png'")
    }


    ti := exif.NewTagIndex()

    it, err := ti.GetWithName(exif.RootIi, "ImageWidth")
    log.PanicIf(err)


    ite := rootIfd.EntriesByTagId[it.Id][0]

    value, err := rootIfd.TagValue(ite)
    log.PanicIf(err)

    expected := []uint32 { 11 }
    if reflect.DeepEqual(value, expected) != true {
        t.Fatalf("image-width not valid")
    }
}

func TestGetExif_Other(t *testing.T) {
    filepath := path.Join(assetsPath, "image.tiff")

    mediaType, rootIfd, err := GetExif(filepath)
    log.PanicIf(err)

    if mediaType != "other" {
        t.Fatalf("media-type not 'other'")
    }


    ti := exif.NewTagIndex()

    it, err := ti.GetWithName(exif.RootIi, "Artist")
    log.PanicIf(err)


    ite := rootIfd.EntriesByTagId[it.Id][0]

    value, err := rootIfd.TagValue(ite)
    log.PanicIf(err)

    expected := "Jean Cornillon"
    if value.(string) != expected {
        t.Fatalf("artist not correct")
    }
}
