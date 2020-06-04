package exifknife

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"io/ioutil"
	"path/filepath"

	"github.com/dsoprea/go-heic-exif-extractor"
	"github.com/dsoprea/go-jpeg-image-structure"
	"github.com/dsoprea/go-png-image-structure"
	"github.com/dsoprea/go-tiff-image-structure"

	"github.com/dsoprea/go-exif/v2"
	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/image"
)

const (
	JpegMediaType  = "jpeg"
	PngMediaType   = "png"
	HeicMediaType  = "heic"
	TiffMediaType  = "tiff"
	OtherMediaType = "other"
)

var (
	ErrNoExif = errors.New("file does not have EXIF")
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

	var jmp riimage.MediaParser
	jmp = jpegstructure.NewJpegMediaParser()

	var pmp riimage.MediaParser
	pmp = pngstructure.NewPngMediaParser()

	var hemp riimage.MediaParser
	hemp = heicexif.NewHeicExifMediaParser()

	var tmp riimage.MediaParser
	tmp = tiffstructure.NewTiffMediaParser()

	mt := ""

	if imageFilepath != "-" {
		extension := filepath.Ext(imageFilepath)
		extension = strings.ToLower(extension)

		if extension == ".jpg" || extension == ".jpeg" {
			mt = JpegMediaType
		} else if extension == ".png" {
			mt = PngMediaType
		} else if extension == ".heic" {
			mt = HeicMediaType
		} else if extension == ".tiff" {
			mt = TiffMediaType
		}
	}

	if mt == "" {
		if jmp.LooksLikeFormat(data) == true {
			mt = JpegMediaType
		} else if pmp.LooksLikeFormat(data) == true {
			mt = PngMediaType
		} else if hemp.LooksLikeFormat(data) == true {
			mt = HeicMediaType
		} else if tmp.LooksLikeFormat(data) == true {
			mt = TiffMediaType
		} else {
			mt = OtherMediaType
		}
	}

	if mt == JpegMediaType {
		mc = &MediaContext{
			MediaType: JpegMediaType,
			RootIfd:   nil,
			RawExif:   nil,
			Media:     nil,
		}

		var parseErr error

		mc.Media, parseErr = jmp.ParseBytes(data)
		if mc.Media == nil && parseErr != nil {
			log.Panic(err)
		}

		rootIfd, rawExif, err := mc.Media.Exif()
		if err != nil {
			// If we had an error before but still got the list of encountered
			// segments back, we would've still checked to see if we had the
			// EXIF data before failing. If we couldn't find or parse it, just
			// panic with the original error.
			if parseErr != nil {
				log.Panic(parseErr)
			}

			if log.Is(err, exif.ErrNoExif) == true {
				return mc, nil
			} else {
				log.Panic(err)
			}
		}

		mc.RootIfd = rootIfd
		mc.RawExif = rawExif
	} else if mt == PngMediaType {
		mc = &MediaContext{
			MediaType: PngMediaType,
			RootIfd:   nil,
			RawExif:   nil,
			Media:     nil,
		}

		mc.Media, err = pmp.ParseBytes(data)
		log.PanicIf(err)

		rootIfd, rawExif, err := mc.Media.Exif()
		if err != nil {
			if log.Is(err, pngstructure.ErrNoExif) == true {
				return mc, nil
			} else {
				log.Panic(err)
			}
		}

		mc.RootIfd = rootIfd
		mc.RawExif = rawExif
	} else if mt == HeicMediaType {
		mc = &MediaContext{
			MediaType: HeicMediaType,
			RootIfd:   nil,
			RawExif:   nil,
			Media:     nil,
		}

		mc.Media, err = hemp.ParseBytes(data)
		log.PanicIf(err)

		rootIfd, rawExif, err := mc.Media.Exif()
		if err != nil {
			if log.Is(err, pngstructure.ErrNoExif) == true {
				return mc, nil
			} else {
				log.Panic(err)
			}
		}

		mc.RootIfd = rootIfd
		mc.RawExif = rawExif
	} else if mt == TiffMediaType {
		mc = &MediaContext{
			MediaType: TiffMediaType,
			RootIfd:   nil,
			RawExif:   nil,
			Media:     nil,
		}

		mc.Media, err = tmp.ParseBytes(data)
		if err != nil {
			if log.Is(err, exif.ErrNoExif) == true {
				return mc, nil
			} else {
				log.Panic(err)
			}
		}

		rootIfd, rawExif, err := mc.Media.Exif()
		log.PanicIf(err)

		mc.RootIfd = rootIfd
		mc.RawExif = rawExif
	} else if mt == OtherMediaType {
		// Brute force.

		rawExif, err := exif.SearchAndExtractExif(data)
		if err != nil {
			if log.Is(err, exif.ErrNoExif) == true {
				log.Panic(ErrNoExif)
			} else {
				log.Panic(err)
			}
		}

		im := exif.NewIfdMappingWithStandard()
		ti := exif.NewTagIndex()

		_, index, err := exif.Collect(im, ti, rawExif)
		log.PanicIf(err)

		mc = &MediaContext{
			MediaType: OtherMediaType,
			RootIfd:   index.RootIfd,
			RawExif:   rawExif,
			Media:     nil,
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

	// TODO(dustin): Add support for adding/updating EXIF content in HEIC files once we find something that mutates HEIC/HEIF data.

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

		err = cs.WriteTo(f)
		log.PanicIf(err)
	} else {
		log.Panicf("media-type not handled for writing")
	}

	return nil
}
