# Tactical EC2 solution to swaggerui

This is a Temporary Fix (TM) until we either move fully to EC2 or support this
functionality under the AWS lambda.

A micro EC2 instance on the free tier was provisioned via the web with an
encrypted version of the private key in this directory.

Access via

$ ssh -i frank-ec2-dev0.pem ubuntu@ec2-18-193-78-190.eu-central-1.compute.amazonaws.com

dp-find-insights-poc-api is running under a detached tmux session under the
ubuntu user (using the environment in ~ubuntu/.direnv)

