package exifknifegps

import (
	"fmt"
	"time"

	"encoding/json"

	"github.com/dsoprea/go-exif"
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-exif-knife"
)

type ExifGps struct {
}

func (eg *ExifGps) ReadGps(imageFilepath string, includeS2Location, printAsJson bool) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	mc, err := exifknife.GetExif(imageFilepath)
	log.PanicIf(err)

	gpsIfd, err := mc.RootIfd.ChildWithIfdPath(exif.IfdPathStandardGps)
	log.PanicIf(err)

	gi, err := gpsIfd.GpsInfo()
	log.PanicIf(err)

	if printAsJson == true {
		distilled := map[string]interface{}{
			"LatitudeDecimal":  gi.Latitude.Decimal(),
			"LongitudeDecimal": gi.Longitude.Decimal(),
			"Altitude":         gi.Altitude,
			"Timestamp":        gi.Timestamp.Format(time.RFC3339),
			"TimestampUnix":    gi.Timestamp.Unix(),
		}

		if includeS2Location == true {
			distilled["S2LocationId"] = gi.S2CellId()
		}

		data, err := json.MarshalIndent(distilled, "", "    ")
		log.PanicIf(err)

		fmt.Println(string(data))
	} else {
		fmt.Printf("%s\n", gi)

		if includeS2Location == true {
			s2LocationId := gi.S2CellId()

			fmt.Printf("\n")
			fmt.Printf("Google S2 Location: [%d]\n", s2LocationId)
		}
	}

	return nil
}
