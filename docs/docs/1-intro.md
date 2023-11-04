---
sidebar_position: 1
---

# Introduction

`golangci-test` aggregates test utilities for Go.
It creates a simple way to run concurrent tests on multiple instances
while using standard Go built-in tools.

:::danger
If you somehow end up on this repo, please wait a few days before
using it. It goes through a lot of changes, and it's not ready for
production use yet.
:::

:::warning
This project is still in early development stage.

**Use it at your own risk**.
:::

## Features

- Split tests into a specified number of groups so that each group:
    - has similar execution time
    - has similar even number of tests in each package
- Run tests based on the generated splits
- Aggregate test results so that e.g. coverage could be generated
- Detect new packages with tests that aren't included in splits
    - either distribute them between the groups and log info to the
      console about the missing packages
    - or fail the test jobs

## Quick start

### Installation

```bash
go install github.com/mflotynski/golangci-test
```

### Usage

To run tests concurrently, you first need to split your tests into
multiple group of packages. The best way to do this would be to run
them and measure how long each tests takes. Finally, tests could be
split into groups of similar execution time.

This is a behaviour that `splitter` module enables.
There are other methods to split tests, but this is supposed to
generate near-optimal splits in reference to the execution time of
test jobs. For other ways to split tests, see the [splitter docs]
page.

The splitter can be run with the following command:

```bash
golangci-test split --type=time --jobs=4 --output=.golangci-test-splits.json
```

This will split all of the tests into 4 groups based on their
execution time.
The results will be saved in `.golangci-test-splits.json` file.

The next step is to run the tests.
This can be done with the following command:

```bash
golangci-test test --instance-index=3 --output=.golangci-test-results.json
```

This will execute the tests from the 3th group (counted from 0 or 1
TODO: FIND WHAT DOES GITLAB RETURNS AS THE FIRST INSTANCE INDEX).

By default, the command will detect test packages that weren't
included in the test splits and distribute them between the groups.