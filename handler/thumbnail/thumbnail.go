package exifknifethumbnail

import (
	"fmt"
	"os"

	"io/ioutil"

	"github.com/dsoprea/go-logging/v2"

	"github.com/dsoprea/go-exif-knife"
)

type ExifThumbnail struct {
}

func (et *ExifThumbnail) writeBytes(outputFilepath string, data []byte) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	if outputFilepath == "-" {
		os.Stdout.Write(data)
	} else {
		err = ioutil.WriteFile(outputFilepath, data, 0644)
		log.PanicIf(err)
	}

	return nil
}

func (et *ExifThumbnail) ExtractThumbnail(imageFilepath, outputFilepath string) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	mc, err := exifknife.GetExif(imageFilepath)
	log.PanicIf(err)

	ifd := mc.RootIfd

	if outputFilepath == "" {
		fmt.Printf("Please provide an output file-path.\n")
		os.Exit(1)
	}

	thumbnailData, err := ifd.NextIfd().Thumbnail()
	log.PanicIf(err)

	err = et.writeBytes(outputFilepath, thumbnailData)
	log.PanicIf(err)

	return nil
}
