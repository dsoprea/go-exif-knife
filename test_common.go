package exifknife

import (
    "os"
    "path"

    "encoding/binary"

    "github.com/dsoprea/go-logging"
)

var (
    TestDefaultByteOrder = binary.BigEndian

    assetsPath = ""
)

func init() {
    goPath := os.Getenv("GOPATH")
    if goPath == "" {
        log.Panicf("GOPATH is empty")
    }

    assetsPath = path.Join(goPath, "src", "github.com", "dsoprea", "go-exif-knife", "assets")
}
