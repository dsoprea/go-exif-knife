module github.com/dsoprea/go-exif-knife

go 1.13

// Development only
// replace github.com/dsoprea/go-png-image-structure => ../go-png-image-structure
// replace github.com/dsoprea/go-jpeg-image-structure => ../go-jpeg-image-structure
// replace github.com/dsoprea/go-utility => ../go-utility

require (
	github.com/dsoprea/go-exif/v2 v2.0.0-20200321225314-640175a69fe4
	github.com/dsoprea/go-jpeg-image-structure v0.0.0-20200322061814-b09a16a9d2b5
	github.com/dsoprea/go-logging v0.0.0-20190624164917-c4f10aab7696
	github.com/dsoprea/go-png-image-structure v0.0.0-20200322062257-e5baf4dcbe12
	github.com/dsoprea/go-utility v0.0.0-20200322055224-4dc0f716e7d0
	github.com/golang/geo v0.0.0-20200319012246-673a6f80352d // indirect
	github.com/jessevdk/go-flags v1.4.0
	golang.org/x/net v0.0.0-20200320220750-118fecf932d8 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)
