package fixture

import (
	"reflect"
	"testing"
)

// FixtureSetup initializes a fixture. This type is exported for experimental
// purposes, e.g., client code can provide an Include implementation. The type
// of the fixture is specified by type parameter T. T must be a pointer type.
//
// Prefer calling [Init] which uses FixtureSetup under the hood.
//
// If T is not a pointer type, no initialization will be performed.
type FixtureSetup[T any] struct {
	// A testing.TB that will be passed to fixtures during initialization
	TB testing.TB
	// The "fixture" to build.
	Fixture T
	// Determines if a value will be managed by fixture. If no `Include` is
	// specified
	Include func(val reflect.Value) bool

	depVals []reflect.Value
	pkgPath string
}

func (f *FixtureSetup[T]) Init() setuper {
	vType := reflect.TypeFor[T]()
	f.TB.Helper()
	if vType.Kind() != reflect.Pointer {
		f.TB.Fatalf("InitFixture: Fixture must be a pointer. Actual type: %s", vType.Name())
	}
	vType = vType.Elem()
	f.pkgPath = vType.PkgPath()
	return f.init(reflect.ValueOf(f.Fixture))
}

func (f *FixtureSetup[T]) init(val reflect.Value) setuper {
	var setups = new(setups)
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return &nullSetuper{}
		}
		val = val.Elem()
	}

	typ := val.Type()
	if typ.Kind() == reflect.Struct {
		setups.append(f.initStruct(val, typ))
	}

	if !val.CanAddr() {
		asAny := val.Interface()
		setups.tryAppend(asAny)
		f.trySetTB(asAny)
	} else {
		// val must be addressable, as both Setup and Init are mutating functions.
		//
		// Val itself may be a non-pointer field in a pointer struct, which
		// means val.Interface() itself is not a Setuper or Initer*, but we can
		// still get a pointer using Addr() because it's inside an addressable
		// struct.
		//
		// \* While a non-pointer val may _implement_ the interfaces, the
		// implementation would be wrong, as they couldn't mutate.

		asAny := val.Addr().Interface()
		setups.tryAppend(asAny)
		f.trySetTB(asAny)
	}
	return setups
}

func (f *FixtureSetup[T]) initStruct(val reflect.Value, typ reflect.Type) *setups {
	var setups = new(setups)
fields:
	for _, field := range reflect.VisibleFields(typ) {
		if len(field.Index) > 1 || !field.IsExported() {
			// Don't set fields of embedded or unexported fields
			continue
		}
		fieldVal := val.FieldByIndex(field.Index)
		if !f.include(fieldVal) {
			continue
		}
		for _, depVal := range f.depVals {
			if field.Type == depVal.Type() {
				fieldVal.Set(depVal)
				continue fields
			}
		}
		if field.Type.Kind() == reflect.Pointer && fieldVal.IsNil() {
			fieldVal.Set(reflect.New(field.Type.Elem()))
			f.depVals = append(f.depVals, fieldVal)
		}
		setups.append(f.init(fieldVal))
	}
	return setups
}

// trySetTB sets the testing.TB instance on the value if it implements
// fixtureInit. If not, the call is ignored
func (s *FixtureSetup[T]) trySetTB(val any) {
	if init, ok := val.(fixtureInit); ok {
		init.SetTB(s.TB)
	}
}
