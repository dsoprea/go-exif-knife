package exifknife

import (
    "os"
    "path"
    "bytes"
    "fmt"

    "encoding/binary"
    "os/exec"

    "github.com/dsoprea/go-logging"
)

var (
    TestDefaultByteOrder = binary.BigEndian

    assetsPath = ""
)

func RunCommand(command_parts ...string) (output []byte, err error) {
    cmd := exec.Command(command_parts[0], command_parts[1:]...)

    b := new(bytes.Buffer)
    cmd.Stdout = b
    cmd.Stderr = b

    err = cmd.Run()
    raw := b.Bytes()

    if err != nil {
        fmt.Printf(string(raw))
        log.Panic(err)
    }

    return b.Bytes(), nil
}

func init() {
    goPath := os.Getenv("GOPATH")
    if goPath == "" {
        log.Panicf("GOPATH is empty")
    }

    assetsPath = path.Join(goPath, "src", "github.com", "dsoprea", "go-exif-knife", "assets")
}
