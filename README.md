# Knight

This is a Go implementation of the [Knight language](https://github.com/knight-lang/knight-lang/).

## Building

Go 1.17 is required.

To build just run `go build` in the root of the project, or run the build script with `go run build.go`.

The build script will run tests, linters, and will inject extra build information that gets printed when `knight version` or `knight -version` is run.

The build script will do things like stripping paths, disabling symbol tables, and disabling DWARF generation by default as well. To include debug information in the build using the build script run with the `-debug` flag:
```sh
go run build.go -debug
```

The `-race` flag can also be used to enable the data race detector:
```sh
go run build.go -debug -race
```

To run tests:
```sh
go run build.go -test
```
...or...
```sh
go run build.go -watch -clear -test
```

You can run `go run build.go -help` to see all of the available options.

### Reckless

You can also build using the `reckless` tag to turn off things like global variable synchronisation.

Synchronisation is on by default and is required for the tests to not generate data races, but when running a normal program in Knight it's actually safe to disable synchronisation due to its single threaded nature.

To enable `reckless` mode build using either...
```sh
go build -tags "reckless"
```
...or...
```sh
go run build.go -tags "reckless"
```

## Usage

Once built you can use the `-e` flag to provide a string expression to run, or you can use the `-f` flag to provide a file to run.

If the `-a` flag is supplied with one of the values `"sexp"`, `"tree"`, or `"waterfall"` then the AST of the supplied program will be printed in that style.

## Profiling

To generate a pprof profile you can run `knight` with either `-p cpu` or `-p mem`.

For example:
```sh
knight -f knight.kn -p cpu
knight -f knight.kn -p mem
```

This will dump profile data into `cpu.pprof` and `mem.pprof` files respectively which can then be inspected using the built-in Go pprof tool, for example:
```sh
go tool pprof -http :8080 cpu.pprof
go tool pprof -http :8080 mem.pprof
```

You can also generate a trace using:
```sh
knight -f knight.kn -p trace
```

This will dump a trace into `trace.out` which can then be inspected using:
```sh
go tool trace -http localhost:8080 trace.out
```
