package lombok

import (
	_ "embed"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

//go:embed testdata/scan_test_1.go
var scanTestCode1 string

func TestScanCode(t *testing.T) {
	var pkgName = "testdata"
	var code = scanTestCode1

	pkg, err := ScanCode(pkgName, code)
	if err != nil {
		t.Errorf("ScanCode() error = %v", err)
		return
	}

	result := dumpPkgInfoForTestScan(pkg)
	expected := `
type=ScanTest
prop=p1, get=P1, set=
prop=p2, get=GetP2, set=
prop=p3, get=GetP3, set=SetP3
`
	if !reflect.DeepEqual(strings.TrimSpace(result), strings.TrimSpace(expected)) {
		t.Errorf("ScanPkgInfo() result = %v, expected %v", result, expected)
	}
}

func dumpPkgInfoForTestScan(pkg *PkgInfo) string {
	var buf strings.Builder
	for _, typ := range pkg.SortedTypes() {
		_, _ = fmt.Fprintf(&buf, "type=%s\n", typ.Name)
		for prop := range typ.Properties() {
			_, _ = fmt.Fprintf(&buf, "prop=%s, get=%s, set=%s\n", prop.Name, prop.Getter, prop.Setter)
		}
	}
	return buf.String()
}
