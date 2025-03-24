# How fixtures work in pytest

_Fixture_ is inspired by fixtured in pytest ([pytest
documentation on fixtures](https://docs.pytest.org/en/6.2.x/fixture.html)) which
separates itself from other test 

Every other structured test framework has some kind of "inheritance" model for
how fixtures can be extended. OOP class based allows sub-classes to extend
fixtures from super classes. [RSpec](), [mocha](), [jasmine]() has a nested
structure, where one group is nested in _one other_ group.

Pytest provides a more flexible mechanism that allows fixtures to describe
different aspects of setting up a test harness, and mix different fixtures in a
much more flexible model.

In pytest, a test can depend on more than one fixture, and fixtures can depend
on other fixtures. The same fixture is only created once, so if a test X depends
on fixture A and B, and fixture A itself depend on fixture B, A's B is the same
as X's B.

In this hypothetical example, a `Processor` write an output file. Some tests may
care about the output file, creatingin initial contents, or verifying the
contents after exercising the SUT. Yet some cases don't care.

```python
@pytest.fixture 
def fake_fs():
    return create_fake_fs()

@pytest.fixture
def output_filename(fake_fs):
    return fake_fs.temp_file_name()

@pytest.fixture 
def processor(fake_fs, output_filename):
    return Processor(fs=fake_fs, file=output_filename)

def test_processor_error_handling(processor):
    # This test doesn't care about 
    result = processor.process(invalid_data)
    assert_result()

def test_process_overwrites_output(processor, output_filename):
    # Put some dummy data in the output file before
    write_test_data(output_filename)
    processor.process()
    assert_output_file(output_filename)
```

The fixture mechanism makes the construction of a fake file system and creating
a unique output file name transparent.

But if the test wants to depend on the output file, it can be added as a
dependency.

_Fixture_ does not try to replicate all behaviour of pytest, just the core idea
how 

Some features in pytest may solve problem that don't exist in Go when you rely
more on actual types, and other problems may have a more natural solution in Go.

E.g., cleanup is implicitly accomplished by integrating to Go's own testing
framework's cleanup functionality.
