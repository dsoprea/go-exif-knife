module github.com/dsoprea/go-exif-knife

go 1.12

// Development only
// replace github.com/dsoprea/go-utility/v2 => ../go-utility/v2
// replace github.com/dsoprea/go-jpeg-image-structure/v2 => ../go-jpeg-image-structure/v2

// replace github.com/dsoprea/go-exif/v3 => ../go-exif/v3
// replace github.com/dsoprea/go-logging/v2 => ../go-logging/v2
// replace github.com/dsoprea/go-exif-extra => ../go-exif-extra

require (
	github.com/dsoprea/go-exif-extra v0.0.0-20210512210440-c683d9263a55
	github.com/dsoprea/go-exif/v3 v3.0.0-20210512043655-120bcdb2a55e
	github.com/dsoprea/go-iptc v0.0.0-20200610044640-bc9ca208b413 // indirect
	github.com/dsoprea/go-jpeg-image-structure/v2 v2.0.0-20210512043942-b434301c6836
	github.com/dsoprea/go-logging v0.0.0-20200710184922-b02d349568dd
	github.com/dsoprea/go-photoshop-info-format v0.0.0-20200610045659-121dd752914d // indirect
	github.com/dsoprea/go-utility/v2 v2.0.0-20200717064901-2fccff4aa15e
	github.com/jessevdk/go-flags v1.4.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
)
