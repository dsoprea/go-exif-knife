package exifkniferead

import (
	"fmt"
	"sort"
	"strings"

	"encoding/json"

	"github.com/dsoprea/go-exif"
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-exif-knife"
)

type ExifRead struct {
}

func (er *ExifRead) Read(imageFilepath string, justTry bool, specificIfdDesignation string, specificTags []string, justPrintValues bool, printAsJson bool) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	mc, err := exifknife.GetExif(imageFilepath)
	log.PanicIf(err)

	if justTry {
		fmt.Printf("%s\n", mc.MediaType)
		return nil
	}

	ifd := mc.RootIfd

	if specificIfdDesignation != "" {
		ifd, err = exif.FindIfdFromRootIfd(ifd, specificIfdDesignation)
		log.PanicIf(err)

		// If we're displaying a particular IFD, don't display any siblings.
		ifd.NextIfd = nil
	}

	included := sort.StringSlice(specificTags)
	included.Sort()

	if printAsJson == true {
		distilled := make(map[string]map[string]interface{})

		err := er.exportIfd(ifd, included, distilled)
		log.PanicIf(err)

		data, err := json.MarshalIndent(distilled, "", "    ")
		log.PanicIf(err)

		fmt.Println(string(data))
	} else {
		if len(included) > 0 {
			ti := exif.NewTagIndex()
			cb := func(ifd *exif.Ifd, tag *exif.IfdTagEntry) error {
				it, err := ti.Get(ifd.Identity(), tag.TagId)

				tagName := ""
				if err == nil {
					tagName = it.Name
				}

				// Unknown tag.
				if tagName == "" {
					return nil
				}

				i := included.Search(tagName)
				if i >= len(included) || included[i] != tagName {
					return nil
				}

				value, err := ifd.TagValue(tag)
				if err != nil {
					if log.Is(err, exif.ErrUnhandledUnknownTypedTag) == true {
						value = "!UNPARSEABLE"
					} else {
						log.Panic(err)
					}
				}

				if justPrintValues == false {
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

				return nil
			}

			err = ifd.EnumerateTagsRecursively(cb)
			log.PanicIf(err)
		} else {
			ifd.PrintTagTree(true)
		}
	}

	return nil
}

func (er *ExifRead) exportIfd(ifd *exif.Ifd, included sort.StringSlice, distilled map[string]map[string]interface{}) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	ti := exif.NewTagIndex()

	for ifd != nil {
		currentIfdDesignation := exif.IfdDesignation(ifd.Ii, ifd.Index)
		currentIfdDesignation = strings.ToLower(currentIfdDesignation)

		for _, tag := range ifd.Entries {
			if tag.ChildIfdName != "" {
				childIfd, err := ifd.ChildWithName(tag.ChildIfdName)
				log.PanicIf(err)

				err = er.exportIfd(childIfd, included, distilled)
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
				ifdMap = map[string]interface{}{
					tagName: value,
				}
			}

			distilled[currentIfdDesignation] = ifdMap
		}

		ifd = ifd.NextIfd
	}

	return nil
}
