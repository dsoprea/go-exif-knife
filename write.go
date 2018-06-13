package exifknife

import (
    "strings"

    "github.com/dsoprea/go-logging"

    "github.com/dsoprea/go-exif"
)

type ExifWrite struct {

}

func (ew *ExifWrite) Write(inputFilepath string, setTagPhrases []string, outputFilepath string) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    mc, err := GetExif(inputFilepath)
    log.PanicIf(err)

    itevr := exif.NewIfdTagEntryValueResolver(mc.RawExif, mc.RootIfd.ByteOrder)
    rootIb := exif.NewIfdBuilderFromExistingChain(mc.RootIfd, itevr)

    ti := exif.NewTagIndex()

    for _, fieldSpec := range setTagPhrases {
        // Split something like "<IFD name>,tag name,value".
        parts := strings.SplitN(fieldSpec, ",", 3)

        ifdDesignation := parts[0]
        tagName := parts[1]
        valueString := parts[2]


        // Validates the IFD designation.
        ini, found := exif.IfdDesignations[ifdDesignation]
        if found == false {
            log.Panicf("IFD designation is not valid: [%s]", ifdDesignation)
        }

        // Validates the tag.
        it, err := ti.GetWithName(ini.Ii, tagName)
        log.PanicIf(err)

        // Ensure we don't have to deal with undefined-type tags at this point in time.
        if it.Type == exif.TypeUndefined {
// TODO(dustin): !! Circle back to this.
            log.Panicf("undefined-type tags are not currently supported for writing")
        }

        tt := exif.NewTagType(it.Type, mc.RootIfd.ByteOrder)

        value, err := tt.FromString(valueString)
        log.PanicIf(err)

        childIb, err := exif.GetOrCreateIbFromRootIb(rootIb, ifdDesignation)

        err = childIb.SetStandardWithName(tagName, value)
        log.PanicIf(err)
    }

    err = SetExif(mc, outputFilepath, rootIb)
    log.PanicIf(err)

    return nil
}
