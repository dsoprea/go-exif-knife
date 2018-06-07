package exifknife

import (
    "strings"

    "path/filepath"

    "github.com/dsoprea/go-jpeg-image-structure"
    "github.com/dsoprea/go-exif"
    "github.com/dsoprea/go-png-image-structure"
    "github.com/dsoprea/go-logging"
)

func GetExif(imageFilepath string) (mediaType string, rootIfd *exif.Ifd, err error) {
    extension := filepath.Ext(imageFilepath)
    extension = strings.ToLower(extension)

    if extension == ".jpg" || extension == ".jpeg" {
        sl, err := jpegstructure.ParseFileStructure(imageFilepath)
        log.PanicIf(err)

        rootIfd, _, err = sl.Exif()
        log.PanicIf(err)

        return "jpeg", rootIfd, nil
    } else if extension == ".png" {
        cs, err := pngstructure.ParseFileStructure(imageFilepath)
        log.PanicIf(err)

        _, rootIfd, err = cs.Exif()
        log.PanicIf(err)

        return "png", rootIfd, nil
    } else {
        rawExif, err := exif.SearchFileAndExtractExif(imageFilepath)
        log.PanicIf(err)

        _, index, err := exif.Collect(rawExif)
        log.PanicIf(err)

        return "other", index.RootIfd, nil
    }
}
