package main

import (
    "os"
    "fmt"
    "strings"
    "sort"
    "time"

    "encoding/json"
    "io/ioutil"

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
    JustValues bool `short:"V" long:"just-values" description:"If specific tags are provided, just print the value for each"`
    Json bool `short:"j" long:"json" description:"Print as JSON"`
}

type writeParameters struct {
    Filepath string `short:"f" long:"filepath" required:"true" description:"File-path ('-' for STDIN)"`
    SetTags []string `short:"s" long:"set-tag" description:"Set tag (can be provided one or more times). Must look like '<ifd:ifd0,ifd1,exif,iop,gps>,<name>,<value>'."`
    OutputFilepath string `short:"o" long:"output-filepath" required:"true" description:"Output file-path ('-' for STDIN)"`
}

type gpsParameters struct {
    Filepath string `short:"f" long:"filepath" required:"true" description:"File-path ('-' for STDIN)"`
    IncludeS2Location bool `short:"2" long:"google-s2" description:"Include Google S2 location"`
    Json bool `short:"j" long:"json" description:"Print as JSON"`
}

type thumbnailParameters struct {
    Filepath string `short:"f" long:"filepath" required:"true" description:"Image file-path ('-' for STDIN)"`
// TODO(dustin): !! When we support updating the thumbnail, try to validate the format.
    OutputFilepath string `short:"o" long:"output-filepath" description:"Thumbnail output file-path ('-' for STDIN)"`
}

type parameters struct {
    Verbose bool `short:"v" long:"verbose" description:"Display logging"`
    Thumbnail thumbnailParameters `command:"thumbnail" alias:"t" description:"Read/write thumbnail"`
    Read readParameters `command:"read" alias:"r" description:"Read/dump EXIF data"`
    Write writeParameters `command:"write" alias:"w" description:"Add/update EXIF data"`
    Gps gpsParameters `command:"gps" alias:"g" description:"Read/dump GPS data from EXIF"`
}

var (
    arguments = new(parameters)
)

func exportIfd(ifd *exif.Ifd, included sort.StringSlice, distilled map[string]map[string]interface{}) (err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    ti := exif.NewTagIndex()

    for ; ifd != nil ; {
        currentIfdDesignation := exif.IfdDesignation(ifd.Ii, ifd.Index)
        currentIfdDesignation = strings.ToLower(currentIfdDesignation)

        for _, tag := range ifd.Entries {
            if tag.ChildIfdName != "" {
                childIfd, err := ifd.ChildWithName(tag.ChildIfdName)
                log.PanicIf(err)

                err = exportIfd(childIfd, included, distilled)
                log.PanicIf(err)

                continue
            }

            it, err := ti.Get(ifd.Identity(), tag.TagId)

            tagName := ""
            if err == nil {
                tagName = it.Name
            }

            // Unknown tag.
            if tagName == "" {
                continue
            }

            i := included.Search(tagName)
            if len(included) > 0 && (i >= len(included) || included[i] != tagName) {
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

            ifdMap, found := distilled[currentIfdDesignation]

            if found == true {
                ifdMap[tagName] = value
            } else {
                ifdMap = map[string]interface{} {
                    tagName: value,
                }
            }

            distilled[currentIfdDesignation] = ifdMap
        }

        ifd = ifd.NextIfd
    }

    return nil
}

func handleRead(options *readParameters) {
    mc, err := exifknife.GetExif(options.Filepath)
    log.PanicIf(err)

    if options.JustTry {
        fmt.Printf("%s\n", mc.MediaType)
        return
    }

    ifd := mc.RootIfd

    if options.SpecificIfd != "" {
        ifd, err = exif.FindIfdFromRootIfd(ifd, options.SpecificIfd)
        log.Panic(err)

        // If we're displaying a particular IFD, don't display any siblings.
        ifd.NextIfd = nil
    }

    included := sort.StringSlice(options.SpecificTags)
    included.Sort()

    if options.Json == true {
        distilled := make(map[string]map[string]interface{})

        err := exportIfd(ifd, included, distilled)
        log.PanicIf(err)

        data, err := json.MarshalIndent(distilled, "", "    ")
        log.PanicIf(err)

        fmt.Println(string(data))
    } else {
        if len(options.SpecificTags) > 0 {
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

                // Unknown tag.
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

                if options.JustValues == false {
                    fmt.Printf("%s: ", tagName)
                }

                switch value.(type) {
                case []uint8:
                    list_ := value.([]uint8)
                    for _, item := range list_ {
                        fmt.Printf("%d ", item)
                    }
                case []uint16:
                    list_ := value.([]uint16)
                    for _, item := range list_ {
                        fmt.Printf("%d ", item)
                    }
                case []uint32:
                    list_ := value.([]uint32)
                    for _, item := range list_ {
                        fmt.Printf("%d ", item)
                    }
                case []int32:
                    list_ := value.([]int32)
                    for _, item := range list_ {
                        fmt.Printf("%d ", item)
                    }
                case []exif.Rational:
                    list_ := value.([]exif.Rational)
                    for _, item := range list_ {
                        fmt.Printf("%d/%d ", item.Numerator, item.Denominator)
                    }
                case []exif.SignedRational:
                    list_ := value.([]exif.SignedRational)
                    for _, item := range list_ {
                        fmt.Printf("%d/%d ", item.Numerator, item.Denominator)
                    }
                case string:
                    fmt.Printf("%s", value.(string))
                default:
                    fmt.Printf("%v", value)
                }

                fmt.Printf("\n")
            }
        } else {
            ifd.PrintTagTree(true)
        }
    }
}

func handleGps(options *gpsParameters) {
    mc, err := exifknife.GetExif(options.Filepath)
    log.PanicIf(err)

    gpsIfd, err := mc.RootIfd.ChildWithIfdIdentity(exif.GpsIi)
    log.PanicIf(err)

    gi, err := gpsIfd.GpsInfo()
    log.PanicIf(err)

    if options.Json == true {
        distilled := map[string]interface{} {
            "LatitudeDecimal": gi.Latitude.Decimal(),
            "LongitudeDecimal": gi.Longitude.Decimal(),
            "Altitude": gi.Altitude,
            "Timestamp": gi.Timestamp.Format(time.RFC3339),
            "TimestampUnix": gi.Timestamp.Unix(),
        }

        if options.IncludeS2Location == true {
            distilled["S2LocationId"] = gi.S2CellId()
        }

        data, err := json.MarshalIndent(distilled, "", "    ")
        log.PanicIf(err)

        fmt.Println(string(data))
    } else {
        fmt.Printf("%s\n", gi)

        if options.IncludeS2Location == true {
            s2LocationId := gi.S2CellId()

            fmt.Printf("\n")
            fmt.Printf("Google S2 Location: [%d]\n", s2LocationId)
        }
    }
}

func writeBytes(outputFilepath string, data []byte) (err error) {
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

func handleThumbnail(options *thumbnailParameters) {
    mc, err := exifknife.GetExif(options.Filepath)
    log.PanicIf(err)

    ifd := mc.RootIfd

    if options.OutputFilepath != "" {
        thumbnailData, err := ifd.NextIfd.Thumbnail()
        log.PanicIf(err)

        err = writeBytes(options.OutputFilepath, thumbnailData)
        log.PanicIf(err)
    } else {
        fmt.Printf("Please provide an output file-path.\n")
        os.Exit(1)
    }
}

func main() {
    p := flags.NewParser(arguments, flags.Default)

    _, err := p.Parse()
    if err != nil {
        os.Exit(1)
    }

    switch p.Active.Name {
    case "read":
        handleRead(&arguments.Read)

    case "write":
        options := &arguments.Write

        ew := new(exifknife.ExifWrite)

        err := ew.Write(options.Filepath, options.SetTags, options.OutputFilepath)
        log.PanicIf(err)

    case "gps":
        handleGps(&arguments.Gps)

    case "thumbnail":
        handleThumbnail(&arguments.Thumbnail)
    }
}
