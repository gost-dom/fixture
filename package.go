// Package fixture helps build "fixture" types for test cases.
//
// Fixture types are test types that help setup the tested code in a
// well-controlled context. This is intended for the cases when
//
// - Multiple tests depend on the same; or similar context.
// - Setting up the tested code in the right context is non-trivial.
//
// During fixture initialization the following will be performed on all "fixture
// types" (see below):
//
// - Iterate all struct fields with exported names.
// - Create new values for nil pointer values.
// - Reuse existing created values if multiple fixtures share the same dependency.
// - Call `SetTB(testing.TB)` on any fixture types that implement it.
// - Aggregate any `Setup()` methods on any fixture type that implement it
//
// Setup functions are not called directly, but a combined Setup method is
// returned to the test, permitting the test to provide additional setup before
// calling the fixture's setup.
//
// Setup is called depth-first, so a fixture can safely assume all fixture types
// it includes have been setup.
//
// The preferred initialization method is to call [Init], but that does not
// permit overriding what is considered a fixture type.
//
// A fixture type is any type that has the "Fixture" suffix. Experimental
// support for overriding this exist in the [FixtureSetup] type directly.
//
// Any SetTB(testing.TB) or Setup() functions **should** be idempotent. When a
// fixture A embeds fixture B, and fixture B implements either SetTB or Setup,
// and A does not, fixture B's method will be called twice. It is not possible
// to change this behaviour in a reliable manner as the reflect package does not
// provide information if an interface is implemented by promoted methods or
// not.
package fixture
