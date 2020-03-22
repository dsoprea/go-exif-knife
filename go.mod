module github.com/dsoprea/go-exif-knife

go 1.13

// Development only
// replace github.com/dsoprea/go-png-image-structure => ../go-png-image-structure
// replace github.com/dsoprea/go-jpeg-image-structure => ../go-jpeg-image-structure

// replace github.com/dsoprea/go-utility => ../go-utility

require (
	github.com/dsoprea/go-exif/v2 v2.0.0-20200321225314-640175a69fe4
	github.com/dsoprea/go-jpeg-image-structure v0.0.0-20200322161747-68b42ab8823c
	github.com/dsoprea/go-logging v0.0.0-20190624164917-c4f10aab7696
	github.com/dsoprea/go-png-image-structure v0.0.0-20200322160151-b1c4e47ed1b3
	github.com/dsoprea/go-utility v0.0.0-20200322154813-27f0b0d142d7
	github.com/golang/geo v0.0.0-20200319012246-673a6f80352d // indirect
	github.com/jessevdk/go-flags v1.4.0
	golang.org/x/net v0.0.0-20200320220750-118fecf932d8 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)
