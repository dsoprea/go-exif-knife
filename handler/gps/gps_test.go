package exifknifegps

import (
	"path"
	"testing"

	"github.com/dsoprea/go-logging/v2"
)

func TestExifGps_ReadGps_WithS2(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "gps.jpg")

	eg := new(ExifGps)

	err := eg.ReadGps(imageFilepath, true, true)
	log.PanicIf(err)
}
