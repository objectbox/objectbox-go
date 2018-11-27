## The [testdata](testdata) folder
All directories under the [testdata](testdata) folder are considered separate test suites and are executed by the `TestAll`
The following rules apply:
* each folder should define entities (a ".go" file)
* each entity should have a special "-binding.expected" file with the expected content of the generated bindings
* each folder should include "objectbox-model.json.expected" - expected content of the ".json" model information
* each folder should include "objectbox-model.go.expected" - expected content of the ".go" model composition code
* there can be a "objectbox-model.json.initial" - it would be used as an initial value for the ".json" file before executing the generator 

### Negative tests
When a file starts with an underscore, it's considered a negative test (the generation should fail):
* it should not generate a binding (thus there's no need for the "-binding.expected" file)


 