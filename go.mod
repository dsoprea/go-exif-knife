module github.com/dsoprea/go-exif-knife

go 1.13

// Development only
// replace github.com/dsoprea/go-utility => ../go-utility
// replace github.com/dsoprea/go-png-image-structure => ../go-png-image-structure
// replace github.com/dsoprea/go-jpeg-image-structure => ../go-jpeg-image-structure
// replace github.com/dsoprea/go-heic-exif-extractor => ../go-heic-exif-extractor
// replace github.com/dsoprea/go-exif/v2 => ../go-exif/v2

require (
	github.com/dsoprea/go-exif/v2 v2.0.0-20200604193436-ca8584a0e1c4
	github.com/dsoprea/go-heic-exif-extractor v0.0.0-20200520190950-3ae4ff88a0d1
	github.com/dsoprea/go-jpeg-image-structure v0.0.0-20200604195911-545332714161
	github.com/dsoprea/go-logging v0.0.0-20200517223158-a10564966e9d
	github.com/dsoprea/go-png-image-structure v0.0.0-20200604195800-fcc57a72b8e7
	github.com/dsoprea/go-tiff-image-structure v0.0.0-20200604200057-42b63f69efa6 // indirect
	github.com/dsoprea/go-utility v0.0.0-20200512094054-1abbbc781176
	github.com/jessevdk/go-flags v1.4.0
)
