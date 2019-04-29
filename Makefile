# Define optional variables here, like GO_TESTARGS.  See
# https://github.com/fastly/gopherup for more info.
BUILDINFO_PKG=github.com/sleach/gstow/pkg/build
include build/gopherup.mk

# Define optional build rules like local-deps and local-clean here.

golangci:
	golangci-lint run

# Define any module-specific targets here
ci: vet golangci test

release:
	goreleaser release --rm-dist
