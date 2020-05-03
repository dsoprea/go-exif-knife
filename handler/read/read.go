package exifkniferead

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"encoding/json"

	"github.com/dsoprea/go-exif/v2"
	"github.com/dsoprea/go-exif/v2/common"
	"github.com/dsoprea/go-logging"

	"github.com/dsoprea/go-exif-knife"
)

var (
	ErrNoExif = errors.New("no EXIF data")
)

type ExifRead struct {
}

func (er *ExifRead) Read(imageFilepath string, justTry bool, specificIfdDesignation string, specificTags []string, justPrintValues bool, printAsJson bool, ignoreNoExif bool) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	s, err := os.Stat(imageFilepath)
	log.PanicIf(err)

	if s.Size() == 0 {
		log.Panicf("zero-length file")
	}

	mc, err := exifknife.GetExif(imageFilepath)
	log.PanicIf(err)

	if mc.RootIfd == nil {
		if ignoreNoExif == true {
			return nil
		}

		return ErrNoExif
	}

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
				it, err := ti.Get(ifd.IfdPath, tag.TagId())

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

				value, err := tag.Value()
				if err != nil {
					if log.Is(err, exifcommon.ErrUnhandledUndefinedTypedTag) == true {
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
				case []exifcommon.Rational:
					list_ := value.([]exifcommon.Rational)
					for _, item := range list_ {
						fmt.Printf("%d/%d ", item.Numerator, item.Denominator)
					}
				case []exifcommon.SignedRational:
					list_ := value.([]exifcommon.SignedRational)
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

	ifdIndex := 0
	for ifd != nil {
		for _, tag := range ifd.Entries {
			if tag.ChildIfdPath() != "" {
				childIfd, err := ifd.ChildWithIfdPath(tag.ChildIfdPath())
				log.PanicIf(err)

				err = er.exportIfd(childIfd, included, distilled)
				log.PanicIf(err)

				continue
			}

			it, err := ti.Get(ifd.IfdPath, tag.TagId())

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

			value, err := tag.Value()
			if err != nil {
				if log.Is(err, exifcommon.ErrUnhandledUndefinedTypedTag) == true {
					value = "!UNPARSEABLE"
				} else {
					log.Panic(err)
				}
			}

			ifdMap, found := distilled[ifd.FqIfdPath]

			if found == true {
				ifdMap[tagName] = value
			} else {
				ifdMap = map[string]interface{}{
					tagName: value,
				}
			}

			distilled[ifd.FqIfdPath] = ifdMap
		}

		ifdIndex++
		ifd = ifd.NextIfd
	}

	return nil
}
