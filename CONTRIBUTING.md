# Contributing

To work on the project:

1. Clone the repository (use recursive clone if you need submodules, e.g. for the client).
2. Install Go and Node per the main [README](README.md) (Building From Source).
3. Run `go mod tidy` to pull Go dependencies.
4. Build the server and admin UI per the README.

To verify changes before submitting:

* Run `golangci-lint run` from the repository root.
* Run `task test unit` to run unit tests
* Run `task test int -a` to run integration tests

For integration tests and other tasks (Docker-based tests, Kubernetes/Tekton), see [test/README.md](test/README.md).
