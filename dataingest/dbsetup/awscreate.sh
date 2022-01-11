#!/bin/bash -x

# gpg -d PGPASSWORD.env.asc for PGPASSWORD
. ../../secrets/PGPASSWORD.env

NAME="fi-database-1"
PORT=54322
REGION="eu-central-1"

# TODO parse from resp
SG="sg-54e0e03d"
HOST="fi-database-1.cbhpmcuqy9vo.eu-central-1.rds.amazonaws.com"

aws --region "$REGION" rds create-db-instance --db-instance-identifier "$NAME" \
--db-instance-class "db.t3.micro"  \
--engine "postgres" --engine-version "13.4" --port "$PORT" \
--master-username "postgres" --master-user-password "$PGPASSWORD"  \
--publicly-accessible --allocated-storage 20 > createlog$$.json

aws --region "$REGION" ec2 authorize-security-group-ingress \
--group-id "$SG" --protocol tcp --port "$PORT" --cidr 0.0.0.0/0
