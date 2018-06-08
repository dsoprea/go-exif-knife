package exifknife

import (
    "fmt"
    "errors"
    // "bytes"
    "time"
    // "strings"
    // "strconv"

    "github.com/dsoprea/go-logging"
    "github.com/dsoprea/go-exif"
)

var (
    ErrNoGpsTags = errors.New("no gps tags")
)


type Degrees struct {
    Orientation byte
    Degrees, Minutes, Seconds int
}

func (d Degrees) String() string {
    return fmt.Sprintf("Degrees<O=[%s] D=(%d) M=(%d) S=(%d)>", string([]byte { d.Orientation }), d.Degrees, d.Minutes, d.Seconds)
}

func (d Degrees) Decimal() float64 {
    decimal := float64(d.Degrees) + float64(d.Minutes) / 60.0 + float64(d.Seconds) / 3600.0

    if d.Orientation == 'S' || d.Orientation == 'W' {
        return -decimal
    } else {
        return decimal
    }
}


type GpsInfo struct {
    Latitude, Longitude Degrees
    Altitude int
    Timestamp time.Time
}

func (gi GpsInfo) String() string {
    return fmt.Sprintf("GpsInfo<LAT=(%.05f) LON=(%.05f) ALT=(%d) TIME=[%s]>", gi.Latitude.Decimal(), gi.Longitude.Decimal(), gi.Altitude, gi.Timestamp)
}

func GetGpsFromExif(gpsIfd *exif.Ifd) (gi *GpsInfo, err error) {
// func GetGpsFromExif(tags []jpegstructure.ExifTag) (gi *GpsInfo, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()


// TODO(dustin): !! Refactor to use the standard go-exif structures.


    return nil, nil

    // indexed := make(map[int]jpegstructure.ExifTag)
    // for _, tag := range tags {
    //     if tag.IfdName != exif.IfdGps || tag.ParentIfdName != exif.IfdStandard {
    //         continue
    //     }

    //     indexed[int(tag.TagId)] = tag
    // }

    // if version, found := indexed[TagVersionId]; found == false {
    //     log.Panic(ErrNoGpsTags)
    // } else if bytes.Compare(version.Value.([]byte), []byte { 2, 2, 0, 0}) != 0 {
    //     log.Panic(ErrNoGpsTags)
    // }

    // latitudeTag, foundLatitude := indexed[TagLatitudeId]
    // latitudeRefTag, foundLatitudeRef := indexed[TagLatitudeRefId]
    // longitudeTag, foundLongitude := indexed[TagLongitudeId]
    // longitudeRefTag, foundLongitudeRef := indexed[TagLongitudeRefId]

    // if foundLatitude != true || foundLatitudeRef != true || foundLongitude != true || foundLongitudeRef != true {
    //     log.Panic(ErrNoGpsTags)
    // }


    // gi = new(GpsInfo)


    // // Parse location.

    // latitudeRaw := latitudeTag.Value.([]exif.Rational)

    // gi.Latitude = Degrees{
    //     Orientation: latitudeRefTag.Value.(string)[0],
    //     Degrees: int(float64(latitudeRaw[0].Numerator) / float64(latitudeRaw[0].Denominator)),
    //     Minutes: int(float64(latitudeRaw[1].Numerator) / float64(latitudeRaw[1].Denominator)),
    //     Seconds: int(float64(latitudeRaw[2].Numerator) / float64(latitudeRaw[2].Denominator)),
    // }

    // longitudeRaw := longitudeTag.Value.([]exif.Rational)

    // gi.Longitude = Degrees{
    //     Orientation: longitudeRefTag.Value.(string)[0],
    //     Degrees: int(float64(longitudeRaw[0].Numerator) / float64(longitudeRaw[0].Denominator)),
    //     Minutes: int(float64(longitudeRaw[1].Numerator) / float64(longitudeRaw[1].Denominator)),
    //     Seconds: int(float64(longitudeRaw[2].Numerator) / float64(longitudeRaw[2].Denominator)),
    // }


    // // Parse altitude.

    // altitudeTag, foundAltitude := indexed[TagAltitudeId]
    // altitudeRefTag, foundAltitudeRef := indexed[TagAltitudeRefId]

    // if foundAltitude == true && foundAltitudeRef == true {
    //     altitudeRaw := altitudeTag.Value.([]exif.Rational)
    //     altitude := int(altitudeRaw[0].Numerator / altitudeRaw[0].Denominator)
    //     if altitudeRefTag.Value.([]byte)[0] == 1 {
    //         altitude *= -1
    //     }

    //     gi.Altitude = altitude
    // }


    // // Parse time.

    // timestampTag, foundTimestamp := indexed[TagTimestampId]
    // datestampTag, foundDatestamp := indexed[TagDatestampId]

    // if foundTimestamp == true && foundDatestamp == true {
    //     dateParts := strings.Split(datestampTag.Value.(string), ":")

    //     year, err1 := strconv.ParseUint(dateParts[0], 10, 16)
    //     month, err2 := strconv.ParseUint(dateParts[1], 10, 8)
    //     day, err3 := strconv.ParseUint(dateParts[2], 10, 8)

    //     if err1 == nil && err2 == nil && err3 == nil {
    //         timestampRaw := timestampTag.Value.([]exif.Rational)

    //         hour := int(timestampRaw[0].Numerator / timestampRaw[0].Denominator)
    //         minute := int(timestampRaw[1].Numerator / timestampRaw[1].Denominator)
    //         second := int(timestampRaw[2].Numerator / timestampRaw[2].Denominator)

    //         gi.Timestamp = time.Date(int(year), time.Month(month), int(day), hour, minute, second, 0, time.UTC)
    //     }
    // }

    // return gi, nil
}
