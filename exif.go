package exifknife

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"io/ioutil"
	"path/filepath"

	"github.com/dsoprea/go-exif/v2"
	"github.com/dsoprea/go-jpeg-image-structure"
	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-png-image-structure"
)

const (
	JpegMediaType  = "jpeg"
	PngMediaType   = "png"
	OtherMediaType = "other"
)

var (
	ErrNoExif = errors.New("file does not have EXIF")
)

// ExifContext is something returned by a MediaParser that knows how to extract
// the actual EXIF structure.
type ExifContext interface {
	// Exif returns the EXIF's root IFD.
	Exif() (rootIfd *exif.Ifd, data []byte, err error)
}

// MediaParser prescribes a specific structure for the parser types that are
// imported from other projects. We don't use it directly, but we use this to
// impose structure.
type MediaParser interface {
	// Parse parses a stream using an `io.Reader`. `ec` should *actually* be a
	// `ExifContext`.
	Parse(r io.Reader, size int) (ec interface{}, err error)

	// ParseFile parses a stream using a file. `ec` should *actually* be a
	// `ExifContext`.
	ParseFile(filepath string) (ec interface{}, err error)

	// ParseBytes parses a stream direct from bytes. `ec` should *actually* be
	// a `ExifContext`.
	ParseBytes(data []byte) (ec interface{}, err error)

	// Parses the data to determine if it's a compatible format.
	LooksLikeFormat(data []byte) bool
}

// MediaContext describes the context/data exteacted from the stream.
type MediaContext struct {
	// MediaType is the name of the detected media type.
	MediaType string

	// RootIfd is the root exif IFD.
	RootIfd *exif.Ifd

	// RawExif is the raw data bytes.
	RawExif []byte

	// Media is type-specific internal data context.
	Media ExifContext
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

	var jmp MediaParser
	jmp = jpegstructure.NewJpegMediaParser()

	var pmp MediaParser
	pmp = pngstructure.NewPngMediaParser()

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
		mc = &MediaContext{
			MediaType: JpegMediaType,
			RootIfd:   nil,
			RawExif:   nil,
			Media:     nil,
		}

		intfc, err := jmp.ParseBytes(data)
		log.PanicIf(err)

		ec := intfc.(ExifContext)
		mc.Media = ec

		rootIfd, rawExif, err := ec.Exif()
		if err != nil {
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

		intfc, err := pmp.ParseBytes(data)
		log.PanicIf(err)

		ec := intfc.(ExifContext)
		mc.Media = ec

		rootIfd, rawExif, err := ec.Exif()
		if err != nil {
			if log.Is(err, pngstructure.ErrNoExif) == true {
				return mc, nil
			} else {
				log.Panic(err)
			}
		}

		mc.RootIfd = rootIfd
		mc.RawExif = rawExif
	} else if mt == OtherMediaType {
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
		log.Panicf("media-type not handled for writing; this shouldn't happen")
	}

	return nil
}
