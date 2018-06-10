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


type MediaContext struct {
    MediaType string
    RootIfd *exif.Ifd
    RawExif []byte
    Media interface{}
}

func GetExif(imageFilepath string) (mc *MediaContext, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    var data []byte

    if imageFilepath == "-" {
        var err error

        data, err = ioutil.ReadAll(os.Stdin)
        log.PanicIf(err)
    } else {
        var err error

        data, err = ioutil.ReadFile(imageFilepath)
        log.PanicIf(err)
    }

// TODO(dustin): !! Add test for all of this.

    jmp := jpegstructure.NewJpegMediaParser()
    pmp := pngstructure.NewPngMediaParser()

    mt := ""
    if imageFilepath != "-" {
        extension := filepath.Ext(imageFilepath)
        extension = strings.ToLower(extension)

        if extension == ".jpg" || extension == ".jpeg" {
            mt = JpegMediaType
        } else if extension == ".png" {
            mt = PngMediaType
        }
    }

    if mt == "" {
        if jmp.LooksLikeFormat(data) == true {
            mt = JpegMediaType
        } else if pmp.LooksLikeFormat(data) == true {
            mt = PngMediaType
        } else {
            mt = OtherMediaType
        }
    }

    if mt == JpegMediaType {
        sl, err := jmp.ParseBytes(data)
        log.PanicIf(err)

        rootIfd, rawExif, err := sl.Exif()
        log.PanicIf(err)

        mc = &MediaContext{
            MediaType: JpegMediaType,
            RootIfd: rootIfd,
            RawExif: rawExif,
            Media: sl,
        }
    } else if mt == PngMediaType {
        cs, err := pmp.ParseBytes(data)
        log.PanicIf(err)

        rootIfd, rawExif, err := cs.Exif()
        log.PanicIf(err)

        mc = &MediaContext{
            MediaType: PngMediaType,
            RootIfd: rootIfd,
            RawExif: rawExif,
            Media: cs,
        }
    } else if mt == OtherMediaType {
        rawExif, err := exif.SearchAndExtractExif(data)
        log.PanicIf(err)

        _, index, err := exif.Collect(rawExif)
        log.PanicIf(err)

        mc = &MediaContext{
            MediaType: OtherMediaType,
            RootIfd: index.RootIfd,
            RawExif: rawExif,
            Media: nil,
        }
    } else {
        log.Panicf("media-type not handled for parsing; this shouldn't happen")
    }

    return mc, nil
}

func SetExif(mc *MediaContext, imageFilepath string, ib *exif.IfdBuilder) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    if mc.Media == nil {
        log.Panicf("media-type does not support writing")
    }

    if mc.MediaType == JpegMediaType {
        sl := mc.Media.(*jpegstructure.SegmentList)

        err := sl.SetExif(ib)
        log.PanicIf(err)

        f, err := os.Create(imageFilepath)
        log.PanicIf(err)

        defer f.Close()

        err = sl.Write(f)
        log.PanicIf(err)
    } else if mc.MediaType == PngMediaType {
        cs := mc.Media.(*pngstructure.ChunkSlice)

        err := cs.SetExif(ib)
        log.PanicIf(err)

        f, err := os.Create(imageFilepath)
        log.PanicIf(err)

        defer f.Close()

        err = cs.Write(f)
        log.PanicIf(err)
    } else {
        log.Panicf("media-type not handled for writing; this shouldn't happen")
    }

    return nil
}
