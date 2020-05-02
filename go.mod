module github.com/dsoprea/go-exif-knife

go 1.13

// Development only
// replace github.com/dsoprea/go-png-image-structure => ../go-png-image-structure
// replace github.com/dsoprea/go-jpeg-image-structure => ../go-jpeg-image-structure
// replace github.com/dsoprea/go-heic-exif-extractor => ../go-heic-exif-extractor
// replace github.com/dsoprea/go-utility => ../go-utility
// replace github.com/dsoprea/go-exif/v2 => ../go-exif/v2

require (
	github.com/dsoprea/go-exif/v2 v2.0.0-20200502211150-cc316fb4407d
	github.com/dsoprea/go-heic-exif-extractor v0.0.0-20200502204916-7d356eae214f
	github.com/dsoprea/go-jpeg-image-structure v0.0.0-20200502204504-6a4a98e0d94b
	github.com/dsoprea/go-logging v0.0.0-20200502201358-170ff607885f
	github.com/dsoprea/go-png-image-structure v0.0.0-20200502204547-95308446cacc
	github.com/dsoprea/go-utility v0.0.0-20200424085841-d6691864fa10
	github.com/jessevdk/go-flags v1.4.0
)
