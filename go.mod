module github.com/dsoprea/go-exif-knife

go 1.13

// Development only
// replace github.com/dsoprea/go-utility/v2 => ../go-utility/v2
// replace github.com/dsoprea/go-png-image-structure/v2 => ../go-png-image-structure/v2
// replace github.com/dsoprea/go-jpeg-image-structure/v2 => ../go-jpeg-image-structure/v2
// replace github.com/dsoprea/go-heic-exif-extractor/v2 => ../go-heic-exif-extractor/v2
// replace github.com/dsoprea/go-tiff-image-structure/v2 => ../go-tiff-image-structure/v2
// replace github.com/dsoprea/go-exif/v3 => ../go-exif/v3
// replace github.com/dsoprea/go-logging/v2 => ../go-logging/v2

require (
	github.com/dsoprea/go-exif/v3 v3.0.0-20200717071058-9393e7afd446
	github.com/dsoprea/go-heic-exif-extractor/v2 v2.0.0-20200717080213-0bb0b9fb3f3a
	github.com/dsoprea/go-iptc v0.0.0-20200610044640-bc9ca208b413 // indirect
	github.com/dsoprea/go-jpeg-image-structure/v2 v2.0.0-20200717072931-d1ef2375db45
	github.com/dsoprea/go-logging v0.0.0-20200710184922-b02d349568dd
	github.com/dsoprea/go-photoshop-info-format v0.0.0-20200610045659-121dd752914d // indirect
	github.com/dsoprea/go-png-image-structure/v2 v2.0.0-20200717073543-83c6344e4591
	github.com/dsoprea/go-tiff-image-structure/v2 v2.0.0-20200717073440-8ac81ec8b423
	github.com/dsoprea/go-utility/v2 v2.0.0-20200717064901-2fccff4aa15e
	github.com/jessevdk/go-flags v1.4.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
)
