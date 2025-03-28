package fixture

import (
	"reflect"
	"strings"
	"testing"
)

// Fixture is a simple default value to embed into custom fixture that need
// access to a [testing.TB] value.
//
//	type MyFixture struct {
//	  fixture.Fixture
//	  *ComponentUnderTest
//	}
//
//	func (f *MyFixture) DataAsStringer() string {
//	  // Data is assumed to always be an fmt.Stringer in this context
//	  if result, ok := f.Data().(fmt.Stringer); ok {
//		return result.String()
//	  }
//	  // embedded testing.TB methods
//	  f.Helper()
//	  f.Error("Data was assume to be a valid fmt.Stringer in this context")
//	  return nil
//	}
//
// Fixture is merely a convience type. Any fixture that implements [FixtureInit]
// will receive the current testing.TB instance.
type Fixture struct{ testing.TB }

func (f *Fixture) SetTB(tb testing.TB) { f.TB = tb }

// FixtureInit should idiomatically have been called setTBer, but that just
// feels awkward.
type FixtureInit interface{ SetTB(testing.TB) }

type setuper interface{ Setup() }

type cleanuper interface{ Cleanup() }

// NullSetuper is just a setuper that ignores setup calls on the null instance
type nullSetuper struct{}

func (s *nullSetuper) Setup()   {}
func (s *nullSetuper) Cleanup() {}

/* -------- setuper -------- */

// setups is a type that lets a slice of setuper's be a setuper themself
type setups []setuper

func (s *setups) append(setup setuper) { *s = append(*s, setup) }

// tryAppend appdns the value to the list of setups if it implements interfaces
// setuper. If not the method does nothing.
func (s *setups) tryAppend(val any) {
	if setup, ok := val.(setuper); ok {
		s.append(setup)
	}
}

func (s setups) Setup() {
	for _, ss := range s {
		ss.Setup()
	}
}

/* -------- cleanuper -------- */

type cleanups []cleanuper

func (c *cleanups) append(cleanup cleanuper) { *c = append(*c, cleanup) }

// tryAppend appends the value to the list of setups if it implements interfaces
// cleanuper. If not the method does nothing.
func (c *cleanups) tryAppend(val any) {
	if cleanup, ok := val.(cleanuper); ok {
		c.append(cleanup)
	}
}

func (c cleanups) Cleanup() {
	for _, cc := range c {
		cc.Cleanup()
	}
}

func (f FixtureSetup[T]) defaultInclude(val reflect.Value) bool {
	var underlyingType = val.Type()
	if underlyingType.Kind() == reflect.Pointer {
		underlyingType = underlyingType.Elem()
	}
	return strings.HasSuffix(underlyingType.Name(), "Fixture")
}

func (f FixtureSetup[T]) include(val reflect.Value) bool {
	if f.Include != nil {
		return f.Include(val)
	}
	return f.defaultInclude(val)
}

// Init initializes the fixture. During initialization it will perform the
// following on all fixture types.
//
// - Call `SetTB(testing.TB)` on any types that implement it.
// - Iterate all struct fields.
// - Create new fixture types
//
// Both `SetTB` and `Setup` functions must be idempotent. When a fixture embeds
// a type that implement one of those methods, the method will be called on both
// the
//
// A fixture type is any type that has the "Fixture" suffix. Experimental
// support for override exist by consuming the [FixtureSetup] type directly.
func Init[T any](
	t testing.TB,
	fixture T,
) (T, setuper) {
	f := &FixtureSetup[T]{
		TB:      t,
		Fixture: fixture,
	}
	return fixture, f.Init()
}
