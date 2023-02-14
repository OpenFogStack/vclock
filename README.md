[![License MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://img.shields.io/badge/License-MIT-brightgreen.svg)
[![Godoc](http://img.shields.io/badge/go-documentation-blue.svg)](https://pkg.go.dev/git.tu-berlin.de/mcc-fred/vclock)

# `vclock` &mdash; A Go Vector Clock Implementation

This repository is a stripped-down version of the brilliant [DistributedClocks/GoVector](github.com/DistributedClocks/GoVector)
library.
Unfortunately, the original library is outdated and has a few bugs that make its use impossible.
This library is created as a drop-in replacement for `github.com/DistributedClocks/GoVector/govec/vclock`.
Other functionality of the package is not implemented.

As of now, these are the major changes:

- include a working `go.mod` file with a specific tag
- fix a bug where `ReturnVCString` is non-deterministic (see [this PR](https://github.com/DistributedClocks/GoVector/pull/67))
- fix a subtle bug in the vector clock comparison function (see [this issue](https://github.com/DistributedClocks/GoVector/issues/68))
- introduce a `Order()` function that returns the relationship of two vector clocks, based on the
  [Voldemort implementation of vector clocks](https://github.com/voldemort/voldemort/blob/master/src/java/voldemort/versioning/VectorClockUtils.java)
- improve documentation

To use this package in your code, download the latest version:

```sh
go get git.tu-berlin.de/mcc-fred/vclock
```

Then replace your import directives from

```go
import (
    "github.com/DistributedClocks/GoVector/govec/vclock"
)
```

to

```go
import (
    "git.tu-berlin.de/mcc-fred/vclock"
)
```

All code continues to be licensed under the [MIT license](./LICENSE).

Note that most development is on the [TU Berlin GitLab instance](https://git.tu-berlin.de/mcc-fred/vclock).
A [GitHub mirror](https://github.com/OpenFogStack/vclock) is provided for convenience.
