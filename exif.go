package exifknife

import (
    "strings"
    "os"

    "path/filepath"
    "io/ioutil"

    "github.com/dsoprea/go-exif"
    "github.com/dsoprea/go-jpeg-image-structure"
    "github.com/dsoprea/go-png-image-structure"
    "github.com/dsoprea/go-logging"
)

const (
    JpegMediaType = "jpeg"
    PngMediaType = "png"
    OtherMediaType = "other"
)

func GetExif(imageFilepath string) (mediaType string, rootIfd *exif.Ifd, err error) {
    if imageFilepath == "-" {

// TODO(dustin): !! Add test for this.

        data, err := ioutil.ReadAll(os.Stdin)
        log.PanicIf(err)

        if jpegstructure.IsJpeg(data) == true {
            sl, err := jpegstructure.ParseBytesStructure(data)
            log.PanicIf(err)

            rootIfd, _, err = sl.Exif()
            log.PanicIf(err)

            return JpegMediaType, rootIfd, nil
        } else if pngstructure.IsPng(data) == true {
            sl, err := pngstructure.ParseBytesStructure(data)
            log.PanicIf(err)

            _, rootIfd, err = sl.Exif()
            log.PanicIf(err)

            return PngMediaType, rootIfd, nil
        } else {
            rawExif, err := exif.SearchAndExtractExif(data)
            log.PanicIf(err)

            _, index, err := exif.Collect(rawExif)
            log.PanicIf(err)

            return OtherMediaType, index.RootIfd, nil
        }
    }

    extension := filepath.Ext(imageFilepath)
    extension = strings.ToLower(extension)

    if extension == ".jpg" || extension == ".jpeg" {
        sl, err := jpegstructure.ParseFileStructure(imageFilepath)
        log.PanicIf(err)

        rootIfd, _, err = sl.Exif()
        log.PanicIf(err)

        return JpegMediaType, rootIfd, nil
    } else if extension == ".png" {
        cs, err := pngstructure.ParseFileStructure(imageFilepath)
        log.PanicIf(err)

        _, rootIfd, err = cs.Exif()
        log.PanicIf(err)

        return PngMediaType, rootIfd, nil
    } else {
        rawExif, err := exif.SearchFileAndExtractExif(imageFilepath)
        log.PanicIf(err)

        _, index, err := exif.Collect(rawExif)
        log.PanicIf(err)

        return OtherMediaType, index.RootIfd, nil
    }
}
