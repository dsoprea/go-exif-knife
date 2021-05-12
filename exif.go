package exifknife

import (
	"fmt"
	"os"

	"io/ioutil"
	"path/filepath"

	"github.com/dsoprea/go-jpeg-image-structure/v2"

	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"github.com/dsoprea/go-logging/v2"
	"github.com/dsoprea/go-utility/v2/image"

	"github.com/dsoprea/go-exif-extra/format"
)

const (
	OtherMediaType = "other"
)

var (
	exifLogger = log.NewLogger("exifknife.exif")
)

// MediaContext describes the context/data exteacted from the stream.
type MediaContext struct {
	// MediaType is the name of the detected media type.
	MediaType string

	// RootIfd is the root exif IFD.
	RootIfd *exif.Ifd

	// RawExif is the raw data bytes.
	RawExif []byte

	// Media is type-specific internal data context.
	Media riimage.MediaContext
}

func (mc MediaContext) String() string {
	hasRootIfd := mc.RootIfd != nil
	return fmt.Sprintf("MediaContext<MEDIA-TYPE=[%s] HAS-ROOT-IFD=[%v] EXIF-LEN=(%d)>", mc.MediaType, hasRootIfd, len(mc.RawExif))
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

	mt := ""
	var mp riimage.MediaParser

	if imageFilepath != "-" {
		extension := filepath.Ext(imageFilepath)

		if extension != "" {
			mt, mp = imageformats.GetFormatForExtension(extension)

			if mp != nil {
				exifLogger.Debugf(nil, "Detected type based on extension: [%s]", mt)
			}
		}
	}

	if mp == nil {
		mt, mp = imageformats.GetFormatForBytes(data)

		if mp != nil {
			exifLogger.Debugf(nil, "Detected type based on content: [%s]", mt)
		}
	}

	if mp != nil {
		mc = &MediaContext{
			MediaType: mt,
			RootIfd:   nil,
			RawExif:   nil,
			Media:     nil,
		}

		mc.Media, err = mp.ParseBytes(data)
		log.PanicIf(err)

		rootIfd, rawExif, err := mc.Media.Exif()
		if err != nil {
			if log.Is(err, exif.ErrNoExif) == true {
				return mc, nil
			} else {
				log.Panic(err)
			}
		}

		mc.RootIfd = rootIfd
		mc.RawExif = rawExif
	} else {
		if mt != "" {
			log.Panicf("Format identified as [%s] but unhandled.", mt)
		}

		// Brute force.

		rawExif, err := exif.SearchAndExtractExif(data)
		log.PanicIf(err)

		im, err := exifcommon.NewIfdMappingWithStandard()
		log.PanicIf(err)

		ti := exif.NewTagIndex()

		_, index, err := exif.Collect(im, ti, rawExif)
		log.PanicIf(err)

		mc = &MediaContext{
			MediaType: OtherMediaType,
			RootIfd:   index.RootIfd,
			RawExif:   rawExif,
			Media:     nil,
		}
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

	// TODO(dustin): Add missing support for adding/updating PNG content.
	// TODO(dustin): Add support for adding/updating EXIF content in HEIC files once we find something that mutates HEIC/HEIF data.

	if mc.MediaType == imageformats.JpegMediaType {
		sl := mc.Media.(*jpegstructure.SegmentList)

		err := sl.SetExif(ib)
		log.PanicIf(err)

		f, err := os.Create(imageFilepath)
		log.PanicIf(err)

		defer f.Close()

		err = sl.Write(f)
		log.PanicIf(err)
	} else {
		log.Panicf("media-type not handled for writing")
	}

	return nil
}
