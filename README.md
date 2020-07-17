[![Build Status](https://travis-ci.org/dsoprea/go-exif-knife.svg?branch=master)](https://travis-ci.org/dsoprea/go-exif-knife)
[![codecov](https://codecov.io/gh/dsoprea/go-exif-knife/branch/master/graph/badge.svg)](https://codecov.io/gh/dsoprea/go-exif-knife)
[![Go Report Card](https://goreportcard.com/badge/github.com/dsoprea/go-exif-knife)](https://goreportcard.com/report/github.com/dsoprea/go-exif-knife)
[![GoDoc](https://godoc.org/github.com/dsoprea/go-exif-knife?status.svg)](https://godoc.org/github.com/dsoprea/go-exif-knife)


## Overview

This is a command-line tool to perform a multitude of surgical operations on the EXIF data in any file that contains it.

This tool has been written on top of [go-exif](https://github.com/dsoprea/go-exif), a complete EXIF implementation.


## Image Support

**JPEG, PNG, HEIC, and TIFF (naturally, since EXIF takes the TIFF structure) are well supported.** A byte-by-byte search will be performed for all other types of file. Writes are only supported for JEPG and PNG.


```
$ ./go-exif-knife write --filepath image.png --set-tag=ifd0,Make,testing --output-filepath /tmp/updated.png

$ ./go-exif-knife read --filepath /tmp/updated.png
 IFD: Ifd<ID=(0) PARENT-IFD=[] IFD=[IFD] INDEX=(0) COUNT=(1) OFF=(0x0008) CHILDREN=(0) PARENT=(0x0000) NEXT-IFD=(0x0000)>
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x010f) TAG-TYPE=[ASCII] UNIT-COUNT=(8)> NAME=[Make] VALUE=[testing]
```


## Install

### Get

- Get via [Github Releases](https://github.com/dsoprea/go-exif-knife/releases)
- Get via go-get and build: "go get github.com/dsoprea/go-exif-knife"

### Install

1. Make the file executable ("chmod 755 [filename]")
2. Put the binary in your path.


## Usage

This tool is comprised of one tool with multiple subcommands. The root tool and the various subcommands provide complete command-line help.


## Examples

### Read

#### Print All Tags

Output (shortened for succinctness):

```
$ ./go-exif-knife read --filepath image.jpg
 IFD: Ifd<ID=(0) PARENT-IFD=[] IFD=[IFD] INDEX=(0) COUNT=(11) OFF=(0x0008) CHILDREN=(2) PARENT=(0x0000) NEXT-IFD=(0x039e)>
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x010f) TAG-TYPE=[ASCII] UNIT-COUNT=(8)> NAME=[Make] VALUE=[samsung]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0110) TAG-TYPE=[ASCII] UNIT-COUNT=(9)> NAME=[Model] VALUE=[SM-N920T]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0112) TAG-TYPE=[SHORT] UNIT-COUNT=(1)> NAME=[Orientation] VALUE=[[1]]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0131) TAG-TYPE=[ASCII] UNIT-COUNT=(12)> NAME=[Software] VALUE=[GIMP 2.8.20]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0132) TAG-TYPE=[ASCII] UNIT-COUNT=(20)> NAME=[DateTime] VALUE=[2018:06:09 01:07:30]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0213) TAG-TYPE=[SHORT] UNIT-COUNT=(1)> NAME=[YCbCrPositioning] VALUE=[[1]]
 - TAG: IfdTagEntry<TAG-IFD=[Exif] TAG-ID=(0x8769) TAG-TYPE=[LONG] UNIT-COUNT=(1)>
   IFD: Ifd<ID=(1) PARENT-IFD=[IFD] IFD=[Exif] INDEX=(0) COUNT=(26) OFF=(0x00d4) CHILDREN=(1) PARENT=(0x0008) NEXT-IFD=(0x0000)>
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x829a) TAG-TYPE=[RATIONAL] UNIT-COUNT=(1)> NAME=[ExposureTime] VALUE=[[{1 13}]]
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x829d) TAG-TYPE=[RATIONAL] UNIT-COUNT=(1)> NAME=[FNumber] VALUE=[[{19 10}]]
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x8822) TAG-TYPE=[SHORT] UNIT-COUNT=(1)> NAME=[ExposureProgram] VALUE=[[2]]
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x9209) TAG-TYPE=[SHORT] UNIT-COUNT=(1)> NAME=[Flash] VALUE=[[0]]
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x920a) TAG-TYPE=[RATIONAL] UNIT-COUNT=(1)> NAME=[FocalLength] VALUE=[[{430 100}]]
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x9286) TAG-TYPE=[UNDEFINED] UNIT-COUNT=(21)> NAME=[UserComment] VALUE=[UserComment<SIZE=(13) ENCODING=[ASCII] V=[0 0 0 73 73 67 83 65]... LEN=(13)>]
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0xa000) TAG-TYPE=[UNDEFINED] UNIT-COUNT=(4)> NAME=[FlashpixVersion] VALUE=[0100]
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0xa001) TAG-TYPE=[SHORT] UNIT-COUNT=(1)> NAME=[ColorSpace] VALUE=[[1]]
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0xa002) TAG-TYPE=[LONG] UNIT-COUNT=(1)> NAME=[PixelXDimension] VALUE=[[920]]
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0xa003) TAG-TYPE=[LONG] UNIT-COUNT=(1)> NAME=[PixelYDimension] VALUE=[[570]]
   - TAG: IfdTagEntry<TAG-IFD=[Iop] TAG-ID=(0xa005) TAG-TYPE=[LONG] UNIT-COUNT=(1)>
     IFD: Ifd<ID=(4) PARENT-IFD=[Exif] IFD=[Iop] INDEX=(0) COUNT=(2) OFF=(0x02b2) CHILDREN=(0) PARENT=(0x00d4) NEXT-IFD=(0x0000)>
     - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0001) TAG-TYPE=[ASCII] UNIT-COUNT=(4)> NAME=[InteroperabilityIndex] VALUE=[R98]
     - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0002) TAG-TYPE=[UNDEFINED] UNIT-COUNT=(4)> NAME=[InteroperabilityVersion] VALUE=[0100]
 - TAG: IfdTagEntry<TAG-IFD=[GPSInfo] TAG-ID=(0x8825) TAG-TYPE=[LONG] UNIT-COUNT=(1)>
   IFD: Ifd<ID=(2) PARENT-IFD=[IFD] IFD=[GPSInfo] INDEX=(0) COUNT=(9) OFF=(0x02d0) CHILDREN=(0) PARENT=(0x0008) NEXT-IFD=(0x0000)>
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0000) TAG-TYPE=[BYTE] UNIT-COUNT=(4)> NAME=[GPSVersionID] VALUE=[[2 2 0 0]]
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0001) TAG-TYPE=[ASCII] UNIT-COUNT=(2)> NAME=[GPSLatitudeRef] VALUE=[N]
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0002) TAG-TYPE=[RATIONAL] UNIT-COUNT=(3)> NAME=[GPSLatitude] VALUE=[[{26 1} {35 1} {12 1}]]
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0006) TAG-TYPE=[RATIONAL] UNIT-COUNT=(1)> NAME=[GPSAltitude] VALUE=[[{0 1}]]
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0007) TAG-TYPE=[RATIONAL] UNIT-COUNT=(3)> NAME=[GPSTimeStamp] VALUE=[[{1 1} {22 1} {57 1}]]
   - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x001d) TAG-TYPE=[ASCII] UNIT-COUNT=(11)> NAME=[GPSDateStamp] VALUE=[2018:04:29]
>IFD: Ifd<ID=(3) PARENT-IFD=[] IFD=[IFD] INDEX=(1) COUNT=(4) OFF=(0x039e) CHILDREN=(0) PARENT=(0x0000) NEXT-IFD=(0x0000)>
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0103) TAG-TYPE=[SHORT] UNIT-COUNT=(1)> NAME=[Compression] VALUE=[[6]]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x011a) TAG-TYPE=[RATIONAL] UNIT-COUNT=(1)> NAME=[XResolution] VALUE=[[{72 1}]]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x011b) TAG-TYPE=[RATIONAL] UNIT-COUNT=(1)> NAME=[YResolution] VALUE=[[{72 1}]]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0128) TAG-TYPE=[SHORT] UNIT-COUNT=(1)> NAME=[ResolutionUnit] VALUE=[[2]]
```


#### Print Specific Tags in All IFDs Using Name

```
$ ./go-exif-knife read --filepath image.jpg --tag Model
Model: SM-N920T
```

"--tag" can be provided multiple times.


#### Print Tag With Name and IFD

```
$ ./go-exif-knife read --filepath "assets/image.jpg" --tag Model --ifd ifd0
Model: Canon EOS 5D Mark III
```

"--tag" can be provided multiple times.


#### Print All Tags in IFD

```
$ ./go-exif-knife read --filepath image.jpg --ifd gps
 IFD: Ifd<ID=(2) PARENT-IFD=[IFD] IFD=[GPSInfo] INDEX=(0) COUNT=(9) OFF=(0x02d0) CHILDREN=(0) PARENT=(0x0008) NEXT-IFD=(0x0000)>
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0000) TAG-TYPE=[BYTE] UNIT-COUNT=(4)> NAME=[GPSVersionID] VALUE=[[2 2 0 0]]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0001) TAG-TYPE=[ASCII] UNIT-COUNT=(2)> NAME=[GPSLatitudeRef] VALUE=[N]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0002) TAG-TYPE=[RATIONAL] UNIT-COUNT=(3)> NAME=[GPSLatitude] VALUE=[[{26 1} {35 1} {12 1}]]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0003) TAG-TYPE=[ASCII] UNIT-COUNT=(2)> NAME=[GPSLongitudeRef] VALUE=[W]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0004) TAG-TYPE=[RATIONAL] UNIT-COUNT=(3)> NAME=[GPSLongitude] VALUE=[[{80 1} {3 1} {13 1}]]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0005) TAG-TYPE=[BYTE] UNIT-COUNT=(1)> NAME=[GPSAltitudeRef] VALUE=[[1]]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0006) TAG-TYPE=[RATIONAL] UNIT-COUNT=(1)> NAME=[GPSAltitude] VALUE=[[{0 1}]]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x0007) TAG-TYPE=[RATIONAL] UNIT-COUNT=(3)> NAME=[GPSTimeStamp] VALUE=[[{1 1} {22 1} {57 1}]]
 - TAG: IfdTagEntry<TAG-IFD=[] TAG-ID=(0x001d) TAG-TYPE=[ASCII] UNIT-COUNT=(11)> NAME=[GPSDateStamp] VALUE=[2018:04:29]
```


#### Print as JSON

```
$ ./go-exif-knife read --filepath image.jpg --ifd gps --json
{
    "gpsinfo": {
        "GPSAltitude": [
            {
                "Numerator": 0,
                "Denominator": 1
            }
        ],
        "GPSAltitudeRef": "AQ==",
        "GPSDateStamp": "2018:04:29",
        "GPSLatitude": [
            {
                "Numerator": 26,
                "Denominator": 1
            },
            {
                "Numerator": 35,
                "Denominator": 1
            },
            {
                "Numerator": 12,
                "Denominator": 1
            }
        ],
        "GPSLatitudeRef": "N",
        "GPSLongitude": [
            {
                "Numerator": 80,
                "Denominator": 1
            },
            {
                "Numerator": 3,
                "Denominator": 1
            },
            {
                "Numerator": 13,
                "Denominator": 1
            }
        ],
        "GPSLongitudeRef": "W",
        "GPSTimeStamp": [
            {
                "Numerator": 1,
                "Denominator": 1
            },
            {
                "Numerator": 22,
                "Denominator": 1
            },
            {
                "Numerator": 57,
                "Denominator": 1
            }
        ],
        "GPSVersionID": "AgIAAA=="
    }
}
```


#### Just Print Value(s)

```
$ ./go-exif-knife read --filepath image.jpg --tag Model --just-values
SM-N920T
```

"--tag" can be provided multiple times.


### GPS

```
$ ./go-exif-knife gps --filepath image.jpg
GpsInfo<LAT=(26.58667) LON=(-80.05361) ALT=(0) TIME=[2018-04-29 01:22:57 +0000 UTC]>
```

```
$ ./go-exif-knife gps --filepath image.jpg --json
{
    "Altitude": 0,
    "LatitudeDecimal": 26.586666666666666,
    "LongitudeDecimal": -80.05361111111111,
    "Timestamp": "2018-04-29T01:22:57Z",
    "TimestampUnix": 1524964977
}
```


Include a geohash calculated with the Google S2 (Hilbert Curve) algorithm:

```
$ ./go-exif-knife gps --filepath image.jpg --json --google-s2
{
    "Altitude": 0,
    "LatitudeDecimal": 26.586666666666666,
    "LongitudeDecimal": -80.05361111111111,
    "S2LocationId": 6542766593732284747,
    "Timestamp": "2018-04-29T01:22:57Z",
    "TimestampUnix": 1524964977
}
```


### Thumbnail

```
$ ./go-exif-knife thumbnail --filepath image.jpg --output-filepath /tmp/thumbnail.jpg
```


### Write

```
$ ./go-exif-knife read --filepath image.jpg --tag Make
Make: samsung

$ ./go-exif-knife write --filepath image.jpg --set-tag=ifd0,Make,testing --output-filepath /tmp/updated.jpg

$ ./go-exif-knife read --filepath /tmp/updated.jpg --tag Make
Make: testing
```
