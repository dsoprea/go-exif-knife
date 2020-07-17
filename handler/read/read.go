package exifkniferead

import (
	"errors"
	"fmt"
	"os"
	"sort"

	"encoding/json"

	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"github.com/dsoprea/go-logging/v2"

	"github.com/dsoprea/go-exif-knife"
)

var (
	readLogger = log.NewLogger("exifkniferead.read")
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

	readLogger.Debugf(nil, "EXIF blob is (%d) bytes.", len(mc.RawExif))

	if justTry {
		fmt.Printf("%s\n", mc.MediaType)
		return nil
	}

	ifd := mc.RootIfd

	if specificIfdDesignation != "" {
		ifd, err = exif.FindIfdFromRootIfd(ifd, specificIfdDesignation)
		log.PanicIf(err)
	}

	included := sort.StringSlice(specificTags)
	included.Sort()

	if printAsJson == true {
		distilled := make(map[string]map[string]interface{})

		err := er.exportIfd(ifd, included, distilled, specificIfdDesignation)
		log.PanicIf(err)

		data, err := json.MarshalIndent(distilled, "", "    ")
		log.PanicIf(err)

		fmt.Println(string(data))
	} else {
		if len(included) > 0 {
			ti := exif.NewTagIndex()
			cb := func(ifd *exif.Ifd, tag *exif.IfdTagEntry) error {
				// This will just add noise to the output (byte-tags are fully
				// dumped).
				if tag.IsThumbnailOffset() == true || tag.IsThumbnailSize() == true {
					return nil
				}

				it, err := ti.Get(ifd.IfdIdentity(), tag.TagId())

				tagName := ""
				if err == nil {
					tagName = it.Name
				}

				i := included.Search(tagName)
				if i >= len(included) || included[i] != tagName {
					return nil
				}

				phrase, err := tag.FormatFirst()
				log.PanicIf(err)

				if justPrintValues == false {
					fmt.Printf("%s: ", tagName)
				}

				fmt.Printf("%s\n", phrase)

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

func (er *ExifRead) exportIfd(ifd *exif.Ifd, included sort.StringSlice, distilled map[string]map[string]interface{}, specificIfdDesignation string) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	ti := exif.NewTagIndex()

	ifdIndex := 0
	for ifd != nil {
		for _, tag := range ifd.Entries() {
			if tag.ChildIfdPath() != "" {
				currentIfdTag := ifd.IfdIdentity().IfdTag()

				childIfdTag :=
					exifcommon.NewIfdTag(
						&currentIfdTag,
						tag.TagId(),
						tag.ChildIfdName())

				iiChild := ifd.IfdIdentity().NewChild(childIfdTag, 0)

				childIfd, err := ifd.ChildWithIfdPath(iiChild)
				log.PanicIf(err)

				err = er.exportIfd(childIfd, included, distilled, specificIfdDesignation)
				log.PanicIf(err)

				continue
			}

			it, err := ti.Get(ifd.IfdIdentity(), tag.TagId())

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

			ifdMap, found := distilled[ifd.IfdIdentity().String()]

			if found == true {
				ifdMap[tagName] = value
			} else {
				ifdMap = map[string]interface{}{
					tagName: value,
				}
			}

			distilled[ifd.IfdIdentity().String()] = ifdMap
		}

		ifdIndex++
		ifd = ifd.NextIfd()

		if specificIfdDesignation != "" {
			// If we're displaying a particular IFD, don't display any siblings.
			break
		}
	}

	return nil
}
