# Unless you set $TEST_TARGET_URL, non-local test targets run against
# the dev EC2 instance, as defined in main.go.

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
	@go test ./... -count=1 -tags=integration

# Usual test (local)
.PHONY: test
test-local:
	@go test ./...  -count=1 -tags=integration -local

# Verbose test
.PHONY: testv
testv:
	 @go test -v ./... -count=1 -tags=integration

# Verbose test (local)
.PHONY: testv
testv-local:
	 @go test -v ./... -count=1 -tags=integration -local


# Very verbose test - displays long coloured diffs for failures
.PHONY: testvv
testvv:
	 @go test -v ./... -count=1 -tags=integration -extra

# Very verbose test - displays long coloured diffs for failures (local)
.PHONY: testvv
testvv-local:
	 @go test -v ./... -count=1 -tags=integration -extra -local

# Benchmark
.PHONY: bench
bench:
	 @go test -bench . -count=1 -tags=integration
