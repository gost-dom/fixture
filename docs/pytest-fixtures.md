# How fixtures work in pytest

The fixture model in pytest is a mini-IOC framework intended for test use cases.
It creates the ability to create common setup code

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

Read more about pytest fixtures on [pytest
documentation](https://docs.pytest.org/en/6.2.x/fixture.html)
