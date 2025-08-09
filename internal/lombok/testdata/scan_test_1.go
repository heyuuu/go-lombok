package testdata

type ScanTest struct {
	p1 string `get:""`
	p2 string `get:"@"`
	p3 string `prop:"@"`
}

func (t ScanTest) P1() string  { return t.p1 }
func (t *ScanTest) P2() string { return t.p2 }
