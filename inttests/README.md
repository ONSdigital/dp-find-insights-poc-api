# inttests

See also Makefile. The following make rules are defined:

### - make update

To add a test add to the Tests literal struct in "main.go".
Also git check-in the new response file under the "resp" directory.

This means we have a git history of API responses as well as allowing diffs.

### - make test & make testv

Runs test and verbose (-v) test. For development work, both have variants to run integration tests against locally-running API server: `make test-local` and `make testv-local`. These assume the API is available at `http://localhost:25252`.

### - make testvv

This is a very verbose test and will display coloured diffs for failures.  This
might be of some use in trying to debug failures with small diffs although the
output is often too long. As above, a local version variant is available for development work, `make testvv-local`, assuming the API is available at `http://localhost:25252`.
### - make bench


This currently runs all tests in sequential and concurrent modes.  There
probably should be more targets and granularity.

## automated tests

The [github actions workflow](https://github.com/features/actions) defined in
`.github/workflows/nightly-integration-tests.yml` performs a nightly run of
`make testvv`
 
