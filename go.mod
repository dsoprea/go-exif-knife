module github.com/dsoprea/go-exif-knife

go 1.13

// Development only
// replace github.com/dsoprea/go-png-image-structure => ../go-png-image-structure
// replace github.com/dsoprea/go-jpeg-image-structure => ../go-jpeg-image-structure
// replace github.com/dsoprea/go-heic-exif-extractor => ../go-heic-exif-extractor
// replace github.com/dsoprea/go-utility => ../go-utility
// replace github.com/dsoprea/go-exif/v2 => ../go-exif/v2

require (
	github.com/dsoprea/go-exif/v2 v2.0.0-20200321225314-640175a69fe4
	github.com/dsoprea/go-heic-exif-extractor v0.0.0-20200323053120-f0dd224ef9dd
	github.com/dsoprea/go-jpeg-image-structure v0.0.0-20200322161747-68b42ab8823c
	github.com/dsoprea/go-logging v0.0.0-20190624164917-c4f10aab7696
	github.com/dsoprea/go-png-image-structure v0.0.0-20200322160151-b1c4e47ed1b3
	github.com/dsoprea/go-utility v0.0.0-20200322184706-df132586647c
	github.com/jessevdk/go-flags v1.4.0
	gopkg.in/yaml.v2 v2.2.8 // indirect
)
