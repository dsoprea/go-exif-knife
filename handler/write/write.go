package exifknifewrite

import (
	"strings"

	"encoding/binary"

	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"github.com/dsoprea/go-logging/v2"

	"github.com/dsoprea/go-exif-knife"
)

type ExifWrite struct {
}

func (ew *ExifWrite) Write(inputFilepath string, setTagPhrases []string, outputFilepath string) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	mc, err := exifknife.GetExif(inputFilepath)
	log.PanicIf(err)

	im, err := exifcommon.NewIfdMappingWithStandard()
	log.PanicIf(err)

	ti := exif.NewTagIndex()

	var rootIb *exif.IfdBuilder

	if mc.RootIfd != nil {
		// There's EXIF data.

		rootIb = exif.NewIfdBuilderFromExistingChain(mc.RootIfd)
	} else {
		// There's no EXIF data. Add it.

		rootIb = exif.NewIfdBuilder(im, ti, exifcommon.IfdStandardIfdIdentity, binary.BigEndian)
	}

	for _, fieldSpec := range setTagPhrases {
		// Split something like "<IFD path>,tag name,value".
		parts := strings.SplitN(fieldSpec, ",", 3)

		ifdPath := parts[0]
		tagName := parts[1]
		valueString := parts[2]

		ii, err := exifcommon.NewIfdIdentityFromString(im, ifdPath)
		log.PanicIf(err)

		// Validates the tag.
		it, err := ti.GetWithName(ii, tagName)
		log.PanicIf(err)

		// Ensure we don't have to deal with undefined-type tags at this point in time.
		if it.DoesSupportType(exifcommon.TypeUndefined) == true {
			// TODO(dustin): !! Circle back to this.
			log.Panicf("undefined-type tags are not currently supported for writing")
		}

		childIb, err := exif.GetOrCreateIbFromRootIb(rootIb, ifdPath)
		log.PanicIf(err)

		err = childIb.SetStandardWithName(tagName, valueString)
		log.PanicIf(err)
	}

	err = exifknife.SetExif(mc, outputFilepath, rootIb)
	log.PanicIf(err)

	return nil
}
