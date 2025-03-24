# Gost-DOM Fixture

_Fixture_ is a tool to help setup test fixtures[^1], i.e., a component in the
test suite that places the SUT[^2] in a controlled context.

This is intended for the case when:

- Setting up the SUT is non-trivial
- Multiple tests needs to use the SUT in an identical or similar manner

_Fixture_ is inspired by [pytest fixtures](./docs/pytest-fixtures.md)

> [!NOTE]
>
> The word "fixture" can be confusing in this document, it can refer to :
>
> - This library (always written in as a proper name, _Fixture_)
> - The default type (always written as a code block `Fixture`).
> - The concept of a fixture (always written in the default typography)
>
> I apologize to screen reader users, that may not easily pick up on the
> typography. I experimented with an alternate name for the library, but I felt
> "fixture" is the right name for the library, and clarity in the readme file
> wasn't a good enough reason to rename.

[^1]: Fixture is a metaphor from mechanical engineering. A fixture holds a piece
    in place, e.g., when testing its mechanical properties. For software testing
    this generally refers to code that places the SUT in a specific context, but
    some test frameworks use the term to refer to the test or suite itself.
[^2]: System under test.

## Looking for sponsors.

This project was conceived as part of the [Gost-DOM](https://gostdom.net)
project, a massively ambitious project to build a headless browser in Go for
testing Go web applications.

Without financial support for the developement, Gost-DOM is unlikely to become a
useful tool.

## How it works.

After reading this, turn to [Getting started](./docs/getting-started.md) for a
code example.

_Fixture_ is based on the following principles.

- A `Fixture` is a component used to setup a _SUT_ in a test context
- A `Fixture` can depend on other fixtures.
- Dependent fixtures can be shared by multiple fixtures.
- A `Fixture` _can_ have initialization/setup code.
- A `Fixture` _can_ integrate to Go's `testing.TB`.

When initializing a fixture, _Fixture_ will iterate through the dependencies, it
will:

- Check if the value should be touched.
- For nil pointer value
  - If a value of this type has already been created, use the same type.
  - Otherwise, create a new empty value
- Unless an configured value was reused:
  - If the type is a struct (or a pointer to one), iterate all fields of the
    struct recursively.
  - If the type has a `SetTB(testing.TB)` method, it will be called.
  - If the type has a `Setup()` method, add this to a list of setup methods.

The `Init` function will return the modified fixture, and a value with a
`Setup()` method, passing control to the test _when_ to setup the method.

## Be aware

> [!DANGER]
>
> _Fixture_ doesn't include a cyclic dependency check. If your fixtures have a
> cyclick dependency, that could possibly result in infinite recursion
> (resulting in a stack overflow)
