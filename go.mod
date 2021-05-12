module github.com/dsoprea/go-exif-knife

go 1.12

// Development only
// replace github.com/dsoprea/go-utility/v2 => ../go-utility/v2
// replace github.com/dsoprea/go-png-image-structure/v2 => ../go-png-image-structure/v2
// replace github.com/dsoprea/go-jpeg-image-structure/v2 => ../go-jpeg-image-structure/v2
// replace github.com/dsoprea/go-heic-exif-extractor/v2 => ../go-heic-exif-extractor/v2
// replace github.com/dsoprea/go-tiff-image-structure/v2 => ../go-tiff-image-structure/v2
// replace github.com/dsoprea/go-webp-image-structure => ../go-webp-image-structure

// replace github.com/dsoprea/go-exif/v3 => ../go-exif/v3
// replace github.com/dsoprea/go-logging/v2 => ../go-logging/v2

require (
	github.com/dsoprea/go-exif/v3 v3.0.0-20210512043655-120bcdb2a55e
	github.com/dsoprea/go-heic-exif-extractor/v2 v2.0.0-20210512044107-62067e44c235
	github.com/dsoprea/go-iptc v0.0.0-20200610044640-bc9ca208b413 // indirect
	github.com/dsoprea/go-jpeg-image-structure/v2 v2.0.0-20210512043942-b434301c6836
	github.com/dsoprea/go-logging v0.0.0-20200710184922-b02d349568dd
	github.com/dsoprea/go-photoshop-info-format v0.0.0-20200610045659-121dd752914d // indirect
	github.com/dsoprea/go-png-image-structure/v2 v2.0.0-20210512044023-23bdd883ee8e
	github.com/dsoprea/go-tiff-image-structure/v2 v2.0.0-20210512044046-dc78da6a809b
	github.com/dsoprea/go-utility/v2 v2.0.0-20200717064901-2fccff4aa15e
	github.com/dsoprea/go-webp-image-structure v0.0.0-20210512044215-f98af2b0401e // indirect
	github.com/jessevdk/go-flags v1.4.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
)
