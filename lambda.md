# Demo Lambda

We've thrown in a quick and dirty lambda for the Find Insights alpha project.
The initial `hello` lambda is a starting point for adding real functionality.

Sources for the lambda are under `functions/hello/`.
Additional lambdas can be created under `functions`/.

The `terraform` directory holds the infrastructure needed to run the lambda.
My initial development uses my ONS AWS account using an S3 bucket for tf
state.
If anybody else wants to run their own lambdas, they will probably need to
set up a different S3 bucket for tf state since bucket names are global.

### Setting up the infra

1. Build the lambda with `make build-lambda bundle-lambda`
2. check the variables make sense in `main.tf`
3. ensure your AWS `development` profile works
4. change to `terraform/` and run `terraform init`
5. run `terraform plan`
5. then run `terraform apply` if it all looks sane

### Lambda development and deployment

For this demo, we're using terraform to deploy the lambda.
Usually S3 is used for deployment artifacts, but I'm using terraform right
now for simplicity.

Dev and testing is basically this:

```
$ vim functions/hello/main.go
$ make build-lambda bundle-lambda
$ ( cd terraform && terraform apply )   # deploys bundle
$ make invoke-lambda                    # runs lambda and prints results
```

You can also test the API gateway pieces like this:

```
$ make invoke-api
$ make invoke-curl
```

The test and deploy cycle might get more involved as we add real endpoints.
