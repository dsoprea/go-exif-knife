package main

import (
	"os"

	"github.com/dsoprea/go-logging/v2"
	"github.com/jessevdk/go-flags"

	"github.com/dsoprea/go-exif-knife/handler/gps"
	"github.com/dsoprea/go-exif-knife/handler/read"
	"github.com/dsoprea/go-exif-knife/handler/thumbnail"
	"github.com/dsoprea/go-exif-knife/handler/write"
)

type readParameters struct {
	Filepath     string   `short:"f" long:"filepath" required:"true" description:"File-path ('-' for STDIN)"`
	JustTry      bool     `short:"s" long:"just-try" description:"Just try to parse. Print file-type and return (0) if success"`
	SpecificIfd  string   `short:"i" long:"ifd" description:"Specific IFD to look at"`
	SpecificTags []string `short:"t" long:"tag" description:"Specific tag to print (can be provided zero or more times)"`
	JustValues   bool     `short:"V" long:"just-values" description:"If specific tags are provided, just print the value for each"`
	Json         bool     `short:"j" long:"json" description:"Print as JSON"`
	IgnoreNoExif bool     `short:"a" long:"ignore-no-exif" description:"Do not print and do not fail if no EXIF data"`
}

type writeParameters struct {
	Filepath string `short:"f" long:"filepath" required:"true" description:"File-path ('-' for STDIN)"`

	// TODO(dustin): !! Provide a way to list the registered IFDs from the command-line and remove them from the help.
	SetTags        []string `short:"s" long:"set-tag" description:"Set tag (can be provided one or more times). Must look like '<ifd path:IFD,IFD0,IFD1,IFD/Exif,IFD/Exif/Iop,IFD/GPSInfo,etc..>,<name>,<value>'."`
	OutputFilepath string   `short:"o" long:"output-filepath" required:"true" description:"Output file-path ('-' for STDIN)"`
}

type gpsParameters struct {
	Filepath          string `short:"f" long:"filepath" required:"true" description:"File-path ('-' for STDIN)"`
	IncludeS2Location bool   `short:"2" long:"google-s2" description:"Include Google S2 location"`
	Json              bool   `short:"j" long:"json" description:"Print as JSON"`
}

type thumbnailParameters struct {
	Filepath string `short:"f" long:"filepath" required:"true" description:"Image file-path ('-' for STDIN)"`
	// TODO(dustin): !! When we support updating the thumbnail, try to validate the format.
	OutputFilepath string `short:"o" long:"output-filepath" description:"Thumbnail output file-path ('-' for STDIN)"`
}

type parameters struct {
	DoPrintVerbose bool `short:"v" long:"verbose" description:"Print logging"`

	Thumbnail thumbnailParameters `command:"thumbnail" alias:"t" description:"Read/write thumbnail"`
	Read      readParameters      `command:"read" alias:"r" description:"Read/dump EXIF data"`
	Write     writeParameters     `command:"write" alias:"w" description:"Add/update EXIF data"`
	Gps       gpsParameters       `command:"gps" alias:"g" description:"Read/dump GPS data from EXIF"`
}

var (
	arguments = new(parameters)
)

func main() {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintError(err)
			os.Exit(1)
		}
	}()

	p := flags.NewParser(arguments, flags.Default)

	_, err := p.Parse()
	if err != nil {
		os.Exit(1)
	}

	if arguments.DoPrintVerbose == true {
		cla := log.NewConsoleLogAdapter()
		log.AddAdapter("console", cla)

		scp := log.NewStaticConfigurationProvider()
		scp.SetLevel(log.LevelDebug)

		log.LoadConfiguration(scp)
	}

	switch p.Active.Name {
	case "read":
		options := &arguments.Read

		er := new(exifkniferead.ExifRead)

		err := er.Read(options.Filepath, options.JustTry, options.SpecificIfd, options.SpecificTags, options.JustValues, options.Json, options.IgnoreNoExif)
		log.PanicIf(err)

	case "write":
		options := &arguments.Write

		ew := new(exifknifewrite.ExifWrite)

		err := ew.Write(options.Filepath, options.SetTags, options.OutputFilepath)
		log.PanicIf(err)

	case "gps":
		options := &arguments.Gps

		eg := new(exifknifegps.ExifGps)

		err := eg.ReadGps(options.Filepath, options.IncludeS2Location, options.Json)
		log.PanicIf(err)

	case "thumbnail":
		options := &arguments.Thumbnail

		et := new(exifknifethumbnail.ExifThumbnail)

		err := et.ExtractThumbnail(options.Filepath, options.OutputFilepath)
		log.PanicIf(err)
	}
}
