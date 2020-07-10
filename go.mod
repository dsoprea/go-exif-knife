module github.com/dsoprea/go-exif-knife

go 1.13

// Development only
// replace github.com/dsoprea/go-utility => ../go-utility
// replace github.com/dsoprea/go-png-image-structure => ../go-png-image-structure
// replace github.com/dsoprea/go-jpeg-image-structure => ../go-jpeg-image-structure
// replace github.com/dsoprea/go-heic-exif-extractor => ../go-heic-exif-extractor
// replace github.com/dsoprea/go-exif/v2 => ../go-exif/v2
// replace github.com/dsoprea/go-logging/v2 => ../go-logging/v2

require (
	github.com/dsoprea/go-exif/v2 v2.0.0-20200710020127-e0ce96b49e0e
	github.com/dsoprea/go-heic-exif-extractor v0.0.0-20200520190950-3ae4ff88a0d1
	github.com/dsoprea/go-iptc v0.0.0-20200610044640-bc9ca208b413 // indirect
	github.com/dsoprea/go-jpeg-image-structure v0.0.0-20200615034914-d40a386309d2
	github.com/dsoprea/go-logging v0.0.0-20200710184922-b02d349568dd
	github.com/dsoprea/go-photoshop-info-format v0.0.0-20200610045659-121dd752914d // indirect
	github.com/dsoprea/go-png-image-structure v0.0.0-20200615034826-4cfc78940228
	github.com/dsoprea/go-tiff-image-structure v0.0.0-20200610044424-2b85e5b2257a
	github.com/dsoprea/go-utility v0.0.0-20200512094054-1abbbc781176
	github.com/go-errors/errors v1.1.1 // indirect
	github.com/jessevdk/go-flags v1.4.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
)
