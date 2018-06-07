package main

import (
    "os"
    "fmt"
    "strings"
    "sort"

    "github.com/jessevdk/go-flags"
    "github.com/dsoprea/go-logging"
    "github.com/dsoprea/go-exif"
    "github.com/dsoprea/go-exif-knife"
)

type readParameters struct {
    Filepath string `short:"f" long:"filepath" required:"true" description:"File-path ('-' for STDIN)"`
    JustTry bool `short:"s" long:"just-try" description:"Just try to parse. Print file-type and return (0) if success"`
    SpecificIfd string `short:"i" long:"ifd" description:"Specific IFD to look at"`
    SpecificTags []string `short:"t" long:"tag" description:"Specific tag to print (can be provided zero or more times)"`
    JustFirstValues bool `short:"v" long:"just-first-values" description:"If specific tags are provided, just print the first value for each (if a list)"`
    JustValues bool `short:"V" long:"just-values" description:"If specific tags are provided, just print the value for each"`
    Json bool `short:"j" long:"json" description:"Print as JSON"`
}

type gpsParameters struct {
}

type parameters struct {
    Verbose bool `short:"v" long:"verbose" description:"Display logging"`
    Read readParameters `command:"read" name:"abc" alias:"r" description:"Read/dump EXIF data"`
    Gps gpsParameters `command:"gps" name:"def" alias:"g" description:"Read/dump GPS data from EXIF"`
}

var (
    arguments = new(parameters)
)

func handleRead() {
    options := arguments.Read

    mediaType, ifd, err := exifknife.GetExif(options.Filepath)
    log.PanicIf(err)

    if options.JustTry {
        fmt.Printf("%s\n", mediaType)
        return
    }

    if options.SpecificIfd != "" {
        ifdName := strings.ToLower(options.SpecificIfd)

        switch ifdName {
        case "ifd0":
        case "exif":
            ifd, err = ifd.ChildWithIfdIdentity(exif.ExifIi)
            log.PanicIf(err)

        case "iop":
            exifIfd, err := ifd.ChildWithIfdIdentity(exif.ExifIi)
            log.PanicIf(err)

            ifd, err = exifIfd.ChildWithIfdIdentity(exif.ExifIopIi)
            log.PanicIf(err)

        case "gps":
            ifd, err = ifd.ChildWithIfdIdentity(exif.GpsIi)
            log.PanicIf(err)

        case "ifd1":
            if ifd.NextIfd == nil {
                log.Panicf("IFD1 not found")
            }

        default:
            fmt.Printf("IFD name not valid. Use 'ifd0', 'ifd1', 'exif', 'gps', or 'iop'.\n")
            os.Exit(2)
        }
    }

    if options.Json == true {

// TODO(dustin): !! Finish.

    } else {
        if len(options.SpecificTags) > 0 {
            included := sort.StringSlice(options.SpecificTags)
            included.Sort()

            ti := exif.NewTagIndex()

            for _, tag := range ifd.Entries {
                // Skip child IFDs. These wouldn't make sense to anyone who
                // doesn't understand EXIF struture.
                if tag.ChildIfdName != "" {
                    continue
                }

                it, err := ti.Get(ifd.Identity(), tag.TagId)

                tagName := ""
                if err == nil {
                    tagName = it.Name
                }

                if tagName == "" {
                    continue
                }

                i := included.Search(tagName)
                if i >= len(included) || included[i] != tagName {
                    continue
                }

                value, err := ifd.TagValue(tag)
                if err != nil {
                    if log.Is(err, exif.ErrUnhandledUnknownTypedTag) == true {
                        value = "!UNPARSEABLE"
                    } else {
                        log.Panic(err)
                    }
                }

// TODO(dustin): !! We should use the same formatting for the JustValues option, too.

                if options.JustFirstValues == true {
                    switch value.(type) {
                    case []uint8:
                        list_ := value.([]uint8)
                        if len(list_) > 0 {
                            fmt.Printf("%d\n", list_[0])
                        }
                    case []uint16:
                        list_ := value.([]uint16)
                        if len(list_) > 0 {
                            fmt.Printf("%d\n", list_[0])
                        }
                    case []uint32:
                        list_ := value.([]uint32)
                        if len(list_) > 0 {
                            fmt.Printf("%d\n", list_[0])
                        }
                    case []int32:
                        list_ := value.([]int32)
                        if len(list_) > 0 {
                            fmt.Printf("%d\n", list_[0])
                        }
                    case []exif.Rational:
                        list_ := value.([]exif.Rational)
                        if len(list_) > 0 {
                            fmt.Printf("%d/%d\n", list_[0].Numerator, list_[0].Denominator)
                        }
                    case []exif.SignedRational:
                        list_ := value.([]exif.SignedRational)
                        if len(list_) > 0 {
                            fmt.Printf("%d/%d\n", list_[0].Numerator, list_[0].Denominator)
                        }
                    case string:
                        fmt.Printf("%s\n", value.(string))
                    default:
                        fmt.Printf("%v\n", value)
                    }
                } else if options.JustValues == true {
                    fmt.Printf("%v\n", value)
                } else {
                    fmt.Printf("%s: %v\n", tagName, value)
                }
            }
        } else {
            ifd.PrintTagTree(true)
        }
    }
}

func handleGps() {


// TODO(dustin): !! Finish.


}

func main() {
    p := flags.NewParser(arguments, flags.Default)

    _, err := p.Parse()
    if err != nil {
        os.Exit(-1)
    }

    switch p.Active.Name {
    case "read":
        handleRead()
    case "gps":
        handleGps()
    }
}
