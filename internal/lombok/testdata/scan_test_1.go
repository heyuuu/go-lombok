package testdata

type ScanTest struct {
	// get / set tag
	p01 string `get:"" set:""`
	p02 string `get:"@" set:""`
	p03 string `get:"AnGetter" set:"AnSetter"`
	// prop tag
	p11 string `prop:""`
	p12 string `prop:"@"`
	p13 string `prop:"name13"`
	p14 string `prop:"@name14"`
	// ref get tag
	p21 string `get:"&"`
	p22 string `get:"&name22"`
	p23 string `prop:"&"`
	p24 string `prop:"&name24"`
}
