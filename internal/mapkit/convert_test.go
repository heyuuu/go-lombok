package mapkit

import "testing"

func TestMapToStruct(t *testing.T) {
	type Anno struct {
		Name   string `name:"name" default:""`
		Age    int    `name:"age"`
		Strict bool   `name:"strict"`
	}

	type args struct {
		m map[string]any
		v any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			"1",
			args{
				m: map[string]any{
					"name":   "lilei",
					"age":    18,
					"strict": true,
				},
				v: &Anno{},
			},
			&Anno{
				Name:   "lilei",
				Age:    18,
				Strict: true,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MapToStruct(tt.args.m, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("MapToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
