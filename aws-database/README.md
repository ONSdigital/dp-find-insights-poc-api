# aws-database

Various scripts to provision & load data into a AWS RDS (Postgres 13.4)
instance used by the Find Insights back-end team.

Most are dependent on the existance of Postgres client utilities being
installed and also a configured aws command line client.

* create.env.asc
  * Encrypted version of postgres password - currently same for "postgres" (admin
user) & "insights" (app user)

Decrypt
```
gpg -d create.env.asc
```

to create "create.env"

* awscreate.sh
  * Tactical solution to create AWS RDS instance and security group opening (non-standard) postgres port of 54322.
  * Probably should be migrated to Terraform.

* awsloaddata.sh
  * creates "insights" pg user & imports "census.sql" DB dump when ran with '-create-user' flag
  * ommitting the flag just imports census.sql"
    * 'dropdb census && createdb census' should be ran first in the latter case

* see READMEDEV.md for details on how to produce a Postgres SQL dump suitable
  for use by 'awsloaddata.sh'
