package exifknifewrite

import (
	"strings"

	"encoding/binary"

	"github.com/dsoprea/go-exif"
	"github.com/dsoprea/go-logging"

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

	var rootIb *exif.IfdBuilder

	// TODO(dustin): !! Why wouldn't this ever be true?
	if mc.RootIfd != nil {
		itevr := exif.NewIfdTagEntryValueResolver(mc.RawExif, mc.RootIfd.ByteOrder)
		rootIb = exif.NewIfdBuilderFromExistingChain(mc.RootIfd, itevr)
	} else {
		im := exif.NewIfdMappingWithStandard()
		ti := exif.NewTagIndex()

		rootIb = exif.NewIfdBuilder(im, ti, exif.IfdPathStandard, binary.BigEndian)
	}

	ti := exif.NewTagIndex()

	for _, fieldSpec := range setTagPhrases {
		// Split something like "<IFD path>,tag name,value".
		parts := strings.SplitN(fieldSpec, ",", 3)

		ifdPath := parts[0]
		tagName := parts[1]
		valueString := parts[2]

		// Validates the tag.
		it, err := ti.GetWithName(ifdPath, tagName)
		log.PanicIf(err)

		// Ensure we don't have to deal with undefined-type tags at this point in time.
		if it.Type == exif.TypeUndefined {
			// TODO(dustin): !! Circle back to this.
			log.Panicf("undefined-type tags are not currently supported for writing")
		}

		tt := exif.NewTagType(it.Type, mc.RootIfd.ByteOrder)

		value, err := tt.FromString(valueString)
		log.PanicIf(err)

		childIb, err := exif.GetOrCreateIbFromRootIb(rootIb, ifdPath)
		log.PanicIf(err)

		err = childIb.SetStandardWithName(tagName, value)
		log.PanicIf(err)
	}

	err = exifknife.SetExif(mc, outputFilepath, rootIb)
	log.PanicIf(err)

	return nil
}
