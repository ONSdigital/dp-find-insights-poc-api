.PHONY: all
all:
	@echo "see README.md or read the Makefile"
	@exit

# Update stored responses under 'resp' directory
.PHONY: update
update:
	@go run -tags=integration .

# Usual test
.PHONY: test
test:
	@go test ./... -tags=integration

# Usual test (local)
.PHONY: test
test-local:
	@go test ./... -tags=integration -local

# Verbose test
.PHONY: testv
testv:
	 @go test -v ./... -tags=integration

# Verbose test (local)
.PHONY: testv
testv-local:
	 @go test -v ./... -tags=integration -local


# Very verbose test - displays long coloured diffs for failures
.PHONY: testvv
testvv:
	 @go test -v ./... -tags=integration -extra

# Very verbose test - displays long coloured diffs for failures (local)
.PHONY: testvv
testvv-local:
	 @go test -v ./... -tags=integration -extra -local

# Benchmark
.PHONY: bench
bench:
	 @go test -bench . -tags=integration
