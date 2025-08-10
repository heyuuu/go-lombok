package lombok

import (
	_ "embed"
	"testing"
)

//go:embed testdata/scan_test_1.go
var scanTestCode1 string

func assertEqual[T comparable](t *testing.T, name string, value T, expected T) bool {
	if value != expected {
		t.Errorf("%s = %v, want %v", name, value, expected)
		return false
	}
	return true
}

func TestScanCode(t *testing.T) {
	var pkgName = "testdata"
	type expectedProperty struct {
		Name        string
		Getter      string
		IsRefGetter bool
		Setter      string
	}
	tests := []struct {
		name          string
		pkg           string
		code          string
		expectedType  string
		expectedProps []expectedProperty
	}{
		{
			pkg:          pkgName,
			code:         scanTestCode1,
			expectedType: "ScanTest",
			expectedProps: []expectedProperty{
				// set / get tag
				{Name: "p01", Getter: "P01", Setter: "SetP01"},
				{Name: "p02", Getter: "GetP02", Setter: "SetP02"},
				{Name: "p03", Getter: "AnGetter", Setter: "AnSetter"},
				//	prop tag
				{Name: "p11", Getter: "P11", Setter: "SetP11"},
				{Name: "p12", Getter: "GetP12", Setter: "SetP12"},
				{Name: "p13", Getter: "Name13", Setter: "SetName13"},
				{Name: "p14", Getter: "GetName14", Setter: "SetName14"},
				// ref tag
				{Name: "p21", Getter: "P21", Setter: "", IsRefGetter: true},
				{Name: "p22", Getter: "name22", Setter: "", IsRefGetter: true},
				{Name: "p23", Getter: "P23", Setter: "SetP23", IsRefGetter: true},
				{Name: "p24", Getter: "Name24", Setter: "SetName24", IsRefGetter: true},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pkg, err := ScanCode(test.pkg, test.code)
			if err != nil {
				t.Errorf("ScanCode(...) error = %v", err)
				return
			}

			typ := pkg.FindType(test.expectedType)
			if typ == nil {
				t.Errorf("FindType(...) = nil, want %v", test.expectedType)
				return
			}

			// 依次检查预期属性
			for _, expectedProp := range test.expectedProps {
				prop := typ.FindProperty(expectedProp.Name)
				if prop == nil {
					t.Errorf("prop `%s` is nil, want not nil", expectedProp.Name)
					continue
				}

				assertEqual(t, "props["+prop.Name+"].Getter", prop.Getter, expectedProp.Getter)
				assertEqual(t, "props["+prop.Name+"].Setter", prop.Setter, expectedProp.Setter)
				assertEqual(t, "props["+prop.Name+"].IsRefGetter", prop.IsRefGetter, expectedProp.IsRefGetter)
			}
		})
	}
}
