module github.com/ONSdigital/dp-find-insights-poc-api

go 1.17

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
	github.com/gosimple/slug v1.12.0
	github.com/jackc/pgx/v4 v4.13.0
	github.com/jtrim-ons/ckmeans v0.0.0-20211215160356-425b5803b027
	github.com/justinas/alice v1.2.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/kylelemons/godebug v1.1.0
	github.com/lib/pq v1.10.3
	github.com/pkg/diff v0.0.0-20210226163009-20ebb0f2a09e
	github.com/pkg/errors v0.9.1
	github.com/smartystreets/goconvey v1.7.2
	github.com/spf13/cast v1.4.1
	github.com/stretchr/testify v1.7.0
	gorm.io/driver/postgres v1.2.1
	gorm.io/gorm v1.22.2
)

require github.com/google/wire v0.5.0

require (
	github.com/ONSdigital/dp-mongodb-in-memory v1.1.0 // indirect
	github.com/cucumber/gherkin-go/v19 v19.0.3 // indirect
	github.com/cucumber/messages-go/v16 v16.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.5 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/gopherjs/gopherjs v0.0.0-20210202160940-bed99a852dfe // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.0 // indirect
	github.com/hashicorp/go-memdb v1.3.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.10.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.1.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.8.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/mailru/easyjson v0.0.0-20190626092158-b2ccc519800e // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/maxcnunes/httpfake v1.2.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.mongodb.org/mongo-driver v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20211124211545-fe61309f8881 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
