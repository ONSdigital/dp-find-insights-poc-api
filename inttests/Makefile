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

# Verbose test
.PHONY: testv
testv:
	 @go test -v ./... -tags=integration

# Very verbose test - displays long coloured diffs for failures
.PHONY: testvv
testvv:
	 @go test -v ./... -tags=integration -args extra


# Benchmark
.PHONY: bench
bench:
	 @go test -bench . -tags=integration
