package main

import (
    "os"
    "fmt"
    "strings"
    "sort"

    "encoding/json"

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
            // We're already on it.

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

            ifd = ifd.NextIfd

        default:
            candidates := make([]string, len(exif.IfdDesignations))
            for key, _ := range exif.IfdDesignations {
                candidates = append(candidates, key)
            }

            fmt.Printf("IFD name not valid. Use: %s\n", strings.Join(candidates, ", "))
            os.Exit(2)
        }

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
