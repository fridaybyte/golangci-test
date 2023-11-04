# Gitlab Integration

There are two ways to integrate.

## Simple approach

Each branch generates the split file during the first pipeline.
Following pipelines use the cached version of the split file.

```yaml
stages:
  - tests

unit_tests:
  stage: tests
  parallel: 2
  script:
    # Gitlab Env variables are detected automatically. That's why
    # there is no need to pass them explicitly by params such as:
    # --split-index=$CI_NODE_INDEX
    # --split-total=$CI_NODE_TOTAL
    - golangci-test test --split-file=.golangci-test-splits.json --fallback=even
  cache:
    policy: pull
  artifacts:
    key: "$CI_COMMIT_REF_NAME"
    paths:
      - test_out_*.json
      - coverage_*.txt

unit_tests_merger:
  stage: tests
  needs: [ "unit_tests" ]
  script:
    # In the following, lines, asterisk is used as a wildcard. 
    # It's a not feature of golangci-test, 
    # but a feature called globbing implemented by most shells. 
    # Read more here: https://tldp.org/LDP/abs/html/globbingref.html 
    - golangci-test split --output=.golangci-test-splits.json --merge ./test_out_*.json
    - golangci-test covmerge ./coverage_*.out > coverage.out
    # TODO: is `out` the standard extension for coverage? :D
  artifacts:
    paths:
      - coverage.txt
  cache:
    policy: push
    paths:
      - .golangci-test-splits.json
```

Pipeline defined in the above *yaml* file will generate pipeline as
shown in the diagram below.

```mermaid
---
title: Flow for the scripts in initial and consecutive pipelines
---
stateDiagram-v2
    U1: Running Gitlab Job "Unit Tests - 1st instance"
    U2: Running Gitlab Job "Unit Tests - 2nd instance"
    M: Merger and the splitter
    state unit_tests_trigger <<fork>>
    state unit_tests_merger <<join>>
    state U1 {
        state if1 <<choice>>
        init1: Trying to load splits file \n from the pulled cache
        load1: Loading splits from the \n splits file pulled from cache
        generate1: Generating test groups \n based on sorted package names
        run1: Running tests from the group 1 and \n storing output in files
        [*] --> init1: Gitlab pulls data from cache
        init1 --> if1
        if1 --> generate1: Splits File \n does not exist
        if1 --> load1: Splits File \n does exist
        generate1 --> run1
        load1 --> run1
        run1 --> [*]: Gitlab stores generated execution time \n and coverage output in artifacts
    }
    state U2 {
        state if2 <<choice>>
        init2: Trying to load splits file \n from the pulled cache
        load2: Loading splits from the \n splits file pulled from cache
        generate2: Generating test groups \n based on sorted package names
        run2: Running tests from the group 2 and \n storing output in files
        [*] --> init2: Gitlab pulls data from cache
        init2 --> if2
        if2 --> generate2: Splits File \n does not exist
        if2 --> load2: Splits File \n does exist
        generate2 --> run2
        load2 --> run2
        run2 --> [*]: Gitlab stores generated execution time \n and coverage output in artifacts
    }
    [*] --> unit_tests_trigger: Start pipeline \n Gitlab triggers unit tests job in parallel
    unit_tests_trigger --> U1
    unit_tests_trigger --> U2
    U1 --> unit_tests_merger
    U2 --> unit_tests_merger
    state M {
        [*] --> init: Gitlab pulls data from artifacts
        init: Trying to load artifacts
        state artifacts_exist <<choice>>
        init --> artifacts_exist
        artifacts_exist --> fail: Data does not exist
        artifacts_exist --> generate: Data exists
        generate: Generating splits and merged coverage \n based on loaded data
        fail: Return error
        generate --> [*]: Gitlab Pushes the generated \n files to the cache
        fail --> [*]: Gitlab Fails the job
    }
    unit_tests_merger --> M: Execution time of tests and \n coverage output is passed in artifacts
    M --> [*]: Finish pipeline

```


