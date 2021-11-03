# stolen from https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-z0-9A-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

LDFLAGS = -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"

.PHONY: all	## run audit, test and build
all: audit test build

.PHONY: audit
audit:	## run nancy auditor
	go list -json -m all | nancy sleuth --exclude-vulnerability-file ./.nancy-ignore

.PHONY: lint
lint:	## doesn't really lint
	exit

.PHONY: build
build:	## build poc service
	go build -tags 'production' $(LDFLAGS) -o $(BINPATH)/dp-find-insights-poc-api

.PHONY: debug
debug:	## run poc service in debug mode
	go build -tags 'debug' $(LDFLAGS) -o $(BINPATH)/dp-find-insights-poc-api
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-find-insights-poc-api

.PHONY: test
test:	## run poc tests
	go test -race -cover ./...

.PHONY: convey
convey:	## run goconvey
	goconvey ./...

.PHONY: test-component
test-component:	## run component tests and coverage
	go test -cover -coverpkg=github.com/ONSdigital/dp-find-insights-poc-api/... -component

#
# these are the lambda-related targets
#

.PHONY: build-lambda
build-lambda:	## compile lambda
	GOOS=linux GOARCH=amd64 go build -o build/hello ./functions/hello/...

.PHONY: bundle-lambda
bundle-lambda:	## bundle lambda into .zip to deploy
	zip -j build/hello.zip build/hello

.PHONY: invoke-lambda
invoke-lambda:	## invoke lambda and display response
	aws --profile development --region eu-central-1 lambda invoke --function-name find-insights-hello .lambda.out
	cat .lambda.out
	rm .lambda.out

.PHONY: invoke-api
invoke-api:	## invoke lambda via api gateway
	REST_API_ID=$$(aws --profile development --region eu-central-1 apigateway get-rest-apis --query 'items[?name==`find-insights-api`]' | jq -r '.[0] .id') ; \
	echo $$REST_API_ID ; \
	RESOURCE_ID=$$(aws --profile development --region eu-central-1 apigateway get-resources --rest-api-id $$REST_API_ID --query 'items[?path==`/hello/{dataset+}`]' | jq -r '.[0] .id') ; \
	echo $$RESOURCE_ID ; \
	aws --profile development --region eu-central-1 \
		apigateway test-invoke-method \
			--rest-api-id $$REST_API_ID \
			--resource-id $$RESOURCE_ID \
			--http-method GET \
			--path-with-query-string /hello/foo

.PHONY: invoke-curl
invoke-curl:	## invoke lambda via curl
	REST_API_ID=$$(aws --profile development --region eu-central-1 apigateway get-rest-apis --query 'items[?name==`find-insights-api`]' | jq -r '.[0] .id') ; \
	echo $$REST_API_ID ; \
	curl --include https://$$REST_API_ID.execute-api.eu-central-1.amazonaws.com/dev/hello/atlas2011.qs101ew
