# Tests

Tests are using [Cucumber](https://cucumber.io/) and [Gherkin syntax](https://cucumber.io/docs/gherkin).

It is using the [Godog](https://github.com/cucumber/godog) framework.

## Run tests

* Authenticate to OCP cluster (need clusteradmin to install crds if not available)
* Run `../hack/run-tests.sh`

### Run specific feature

`../hack/run-tests.sh --feature {FEATURE}`

### Run in concurrent mode

`../hack/run-tests.sh --concurrent {NB_OF_CONCURRENT_TESTS}`

### Other options

Check the usage of the tests script:

`../hack/run-tests.sh -h`
