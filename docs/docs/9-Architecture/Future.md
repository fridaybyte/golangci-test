# Future

## What's next?

### Highest Priority

First of all, `golangci-test` intends to be stable and reliable tool.
This is why the highest priority will always be catching and fixing
any bugs that would stand in a way of including `golangci-test` in
production pipelines.

New features will always have lower priority and their implementation
should be preceded with careful reasoning and thinking of all edge
cases.

### Eventual refactor of merge tools

[Go 1.20 introduced](https://go.dev/blog/integration-test-coverage) a
few great tools focused on integration testing.
Including a tool, that enables merging coverage profiles.
Unfortunately, it supports only new coverage format,
while there is no way (production-ready) way to generate it for
standard unit tests without building test executable and using
unexported (unofficial) flags,
[read more here](https://github.com/golang/go/issues/51430#issuecomment-1344711300).

Development of these features is definitely worth observing as it
could replace custom code with the built-in tools (officially
supported tools)