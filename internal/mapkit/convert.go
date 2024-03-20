package mapkit

import (
	"fmt"
	"reflect"
)

func MapToStruct(m map[string]any, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("cannot convert %T type", v)
	}

	structValue := rv.Elem()
	typePropSetters := getStructFieldSetters(structValue.Type())
	if typePropSetters == nil {
		return nil
	}

	for key, val := range m {
		typePropSetter := typePropSetters[key]
		if typePropSetter == nil {
			continue
		}
		typePropSetter(structValue, val)
	}

	return nil
}

type fieldSetter = func(s reflect.Value, value any)
type fieldSetters = map[string]fieldSetter

var structFieldSettersCache = map[reflect.Type]fieldSetters{}

func getStructFieldSetters(typ reflect.Type) fieldSetters {
	if typ.Kind() != reflect.Struct {
		return nil
	}

	if structFieldSettersCache[typ] != nil {
		return structFieldSettersCache[typ]
	}

	setters := make(fieldSetters)
	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)
		name := structField.Tag.Get("name")
		setters[name] = func(s reflect.Value, val any) {
			valType := reflect.TypeOf(val)
			if valType.ConvertibleTo(structField.Type) {
				newVal := reflect.ValueOf(val).Convert(structField.Type)
				s.FieldByIndex(structField.Index).Set(newVal)
			}
		}
	}

	structFieldSettersCache[typ] = setters
	return structFieldSettersCache[typ]
}
