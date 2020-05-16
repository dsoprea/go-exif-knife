module github.com/dsoprea/go-exif-knife

go 1.13

// Development only
// replace github.com/dsoprea/go-png-image-structure => ../go-png-image-structure
// replace github.com/dsoprea/go-jpeg-image-structure => ../go-jpeg-image-structure
// replace github.com/dsoprea/go-heic-exif-extractor => ../go-heic-exif-extractor
// replace github.com/dsoprea/go-utility => ../go-utility
// replace github.com/dsoprea/go-exif/v2 => ../go-exif/v2

require (
	github.com/dsoprea/go-exif/v2 v2.0.0-20200516122116-a45cc7cfd55e
	github.com/dsoprea/go-heic-exif-extractor v0.0.0-20200502204916-7d356eae214f
	github.com/dsoprea/go-jpeg-image-structure v0.0.0-20200516044922-d03f8ff1b584
	github.com/dsoprea/go-logging v0.0.0-20200502201358-170ff607885f
	github.com/dsoprea/go-png-image-structure v0.0.0-20200502204547-95308446cacc
	github.com/dsoprea/go-utility v0.0.0-20200512094054-1abbbc781176
	github.com/jessevdk/go-flags v1.4.0
	golang.org/x/net v0.0.0-20200513185701-a91f0712d120 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
