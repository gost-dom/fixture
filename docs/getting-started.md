# Getting started with Gost-DOM Fixture

To use _Fixture_, create different fixture types, each must have the "`Fixture`"
suffix (this behaviour is planned to be customizable).

This example tests a component that depends on a file system. Fixtures are
created to setup the component with a fake implementation of `fs.FS`.

```go
type FakeFSFixture struct {
    fstest.MapFS
}

type MyFixture struct {
    *FakeFSFixture
    *MyComponent
}
```

Because `FaksFSFixture` _embeds_ `fstest.MapFS`, a valid `fs.FS` implementation,
`FakeFSFixture` is itself a valid `fs.FS` implementation.


The fake file system must be created _before_ creating the component, so add an
idempotent `Setup()` function to the file system fixture (for this case,
idempotency actually becomes useful later)

```
func (f *FakeFSFixture) Setup() {
    if f.FS == nil {
        f.FS = new(fstest.MapFS)
    }
}
```

It can be useful to automatically create the component under test, so we add a
`Setup()` for `MyFixture` as well:

```go
func (f *MyFixture) Setup() {
    f.MyComponent = &MyComponent{FS: f.FakeFSFixture}
}
```

_Fixture_ guarantees that dependencies are setup before their dependees, so the
fixture can safely assume the `FakeFSFixture` is valid.

Now we can write a test that uses the fixture. _Fixture_ doesn't run call any
`Setup()` functions yet, but let test code control when to call the setup.

```go
func TestMyComponent(t *testing.T) {
    fixture, ctrl := fixture.Init(t, &MyFixture{})
    ctrl.Setup()

    // Process a non-existing file
    res, err := fixture.CreateFile(createTmpFilename())
    // verify ...
}
```

The test controlled setup allows test code to run additional setup _before_
running setup. In this test, the test creates a specific test file. To help
simplify setting up the fake file system, a new helper is added to
`FakeFSFixture`. Now the idempotency of `Setup()` becomes useful:

```go
func (f *FakeFSFixture) CreateFile(path, data string) {
    f.Setup()
    f.FS[path] = fstest.MapFile{ Data: []byte(data) }
}

func TestMyComponent(t *testing.T) {
    fixture, ctrl := fixture.Init(t, &MyFixture{})
    filename := createTempFilename())
    fixture.CreateFile(filename, "Lorem ipsom")

    ctrl.Setup()
    fixture.ProcessFile(filename)
    file := fixture.Open(filename)
    // Verify state of file
}
```

## Adhoc components

The previous examples used embedding on all nested Fixtures, so the fixture
itself implemented all methods of both `fs.FS` and `*MyComponent`. This may not
be a sensible approach when the fixtures grow, and named fields may be a more
sensible approac.

As pointer values to fixtures are reused, any piece of a fixture can just depend
directly on another fixture type to avoid long dot-chains like
`f.XFixture.YFixture.ZFixture.Do()`.

In the test code itself, you can define an inline struct type, to specifically
determine the dependencies used

```go
func TestOutputFile(t *testing.T) {
    f, ctrl := fixture.Init(t, &struct{
        MyComponentFixture
        fs: *FakeFSFixture // Same Fake FS as the component uses
    }{})
    // Allow additional setup before the code's setup is called.
    f.fs.CreateFile("temp.txt", "Lorem ipsum ...")
    ctrl.Setup()
}
```

Support for this pattern is the primary reason `Init` returns two values.

## `testing` integration.

A fixture may contain helper method that assume a specific context; and an
invalid assumption isn't part of the test itself. If an invalid assumption is
encountered, the error should be reported back to the `*testing.T`/`*testing.B`
instance.

The easy solution is to embed `fixture.Fixture` into your own fixture component.
This in turn embeds, `testing.TB`.

```go
type MyFixture struct {
    fixture.Fixture // New! Embed fixture.Fixture
    *FakeFSFixture
    *MyComponent
}

func (f *MyFixture) GetValueThatIsAssumedToBeAFoo() Foo {
    foo, ok := f.MyComponent.Value().(Foo)
    if ok { 
        return foo
    }
    f.Helper()
    f.Errorf("Error getting Foo: Value() did not return a Foo")
    return nil
}
```

_Fixture_ checks of a component implements `FixtureInit`. `fixture.Fixture` is
merely an easy default implementation

```Go
type FixtureInit interface {
    SetTB(testing.TB)
}
```

## Notes

> [!DANGER]
>
> `Setup()` and `SetTB(testing.TB)` functions _may be called twice_ so it's
> imperative that they're idempotent.


- "idempotent"/ "idempotency" - An operation is said to be idempotent if
re-running the same operations with the same arguments does not
produce any new side effect.
