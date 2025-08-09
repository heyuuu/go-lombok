package lombok

import (
	_ "embed"
	"testing"
)

//go:embed testdata/test_1.go
var genTest1Code string

//go:embed testdata/test_1.properties.go
var genTest1Expected string

func TestGenerateByCode(t *testing.T) {
	var pkgName = "testdata"
	var code = genTest1Code
	var expected = genTest1Expected

	result, err := GenerateByCode(pkgName, code)
	if err != nil {
		t.Errorf("GenerateByCode() error = %v", err)
		return
	}

	if result != expected {
		t.Errorf("GenerateByCode() = %v, want %v", result, expected)
	}
}
