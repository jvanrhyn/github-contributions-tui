# The Go binary to use for building and running commands
BIN=go

# The output path for the build artifacts
OUTPATH=./bin

# The build target that compiles the Go code and outputs the binary to the specified OUTPATH
build: create_build_folder copy_env
    ${BIN} build -v -o ${OUTPATH} ./...

# The target to create the build folder if it does not exist
create_build_folder:
    mkdir -p bin

# The target to copy the .env file to the output path if it exists
copy_env:
    @envfile=$$(find . -name ".env" -print -quit); \
    if [ -f "$${envfile}" ] && [ -f "${OUTPATH}/.env" ]; then \
        rm -f "${OUTPATH}/.env" || exit 1; \
    fi; \
    if [ -f "$${envfile}" ]; then \
        cp -f "$${envfile}" "${OUTPATH}/" || exit 1; \
    fi

# The target to run tests with race detection and verbose output
test:
    go test -race -v ./...

# The target to watch for changes and run tests automatically
watch-test:
    reflex -t 50ms -s -- sh -c 'gotest -race -v ./...'

# The target to run benchmark tests with memory allocation statistics
bench:
    go test -benchmem -count 3 -bench ./...

# The target to watch for changes and run benchmark tests automatically
watch-bench:
    reflex -t 50ms -s -- sh -c 'go test -benchmem -count 3 -bench ./...'

# The target to generate a test coverage report and output it as an HTML file
coverage:
    ${BIN} test -v ./... -coverprofile=cover.out -covermode=atomic
    ${BIN} tool cover -html=cover.out -o cover.html

# The target to install various development tools
tools:
    ${BIN} install github.com/cespare/reflex@latest
    ${BIN} install github.com/rakyll/gotest@latest
    ${BIN} install github.com/psampaz/go-mod-outdated@latest
    ${BIN} install github.com/jondot/goweight@latest
    ${BIN} install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    ${BIN} get -t -u golang.org/x/tools/cmd/cover
    ${BIN} install github.com/sonatype-nexus-community/nancy@latest
    go mod tidy

# The target to run the linter with a timeout and limit on the number of same issues
lint:
    golangci-lint run --timeout 60s --max-same-issues 50 ./...

# The target to run the linter and automatically fix issues
lint-fix:
    golangci-lint run --timeout 60s --max-same-issues 50 --fix ./...

# The target to audit dependencies for vulnerabilities
audit:
    ${BIN} list -json -m all | nancy sleuth

# The target to list outdated dependencies and suggest updates
outdated:
    ${BIN} list -u -m -json all | go-mod-outdated -update -direct

# The target to analyze the size of the Go binary
weight:
    goweight