module github.com/ONSdigital/dp-find-insights-poc-api

go 1.16

replace github.com/coreos/etcd => github.com/coreos/etcd v3.3.24+incompatible

require (
	github.com/ONSdigital/dp-api-clients-go v1.41.1
	github.com/ONSdigital/dp-component-test v0.6.0
	github.com/ONSdigital/dp-healthcheck v1.1.3
	github.com/ONSdigital/dp-net v1.2.0
	github.com/ONSdigital/log.go/v2 v2.0.9
	github.com/aws/aws-lambda-go v1.27.0
	github.com/aws/aws-sdk-go v1.41.15
	github.com/cucumber/godog v0.12.1
	github.com/deepmap/oapi-codegen v1.9.0
	github.com/getkin/kin-openapi v0.86.0
	github.com/go-chi/chi/v5 v5.0.4
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/jackc/pgx/v4 v4.13.0
	github.com/justinas/alice v1.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/klauspost/compress v1.12.3 // indirect
	github.com/kylelemons/godebug v1.1.0
	github.com/lib/pq v1.10.3
	github.com/pkg/diff v0.0.0-20210226163009-20ebb0f2a09e
	github.com/pkg/errors v0.9.1
	github.com/smartystreets/goconvey v1.6.6
	github.com/spf13/cast v1.4.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/net v0.0.0-20210913180222-943fd674d43e // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	gorm.io/driver/postgres v1.2.1
	gorm.io/gorm v1.22.2
)
