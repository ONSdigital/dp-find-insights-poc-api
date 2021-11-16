# inttests

See also Makefile

## make update

To add a test add to the Tests literal struct in "main.go" (including the correct
SHA1 hash for the expected response).  Also git check-in the new response file
under the "resp" directory.

This means we have a git history of API responses as well as allowing diffs.

## make test & make testv

Runs test and verbose (-v) test.

## make testvv

This is a very verbose test and will display coloured diffs for failures.  This
might be of some use in trying to debug failures with small diffs although the
output is often too long.

## make bench

This currently runs all tests in sequential and concurrent modes.  There
probably should be more targets and granularity.

