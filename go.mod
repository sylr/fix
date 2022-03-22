module sylr.dev/fix

go 1.18

require (
	filippo.io/age v1.0.0
	github.com/google/uuid v1.3.0
	github.com/iancoleman/strcase v0.2.0
	github.com/mattn/go-sqlite3 v1.14.12
	github.com/nats-io/nats-server/v2 v2.7.4
	github.com/nats-io/nats.go v1.13.1-0.20220308171302-2f2f6968e98d
	github.com/quickfixgo/enum v0.0.0-20210629025633-9afc8539baba
	github.com/quickfixgo/field v0.0.0-20171007195410-74cea5ec78c7
	github.com/quickfixgo/fix50sp1 v0.0.0-20171007213322-66c7df13880d
	github.com/quickfixgo/fix50sp2 v0.0.0-20171007213400-aacf4747f57a
	github.com/quickfixgo/quickfix v0.6.1-0.20210618140103-31f5ebe90229
	github.com/quickfixgo/tag v0.0.0-20171007194743-cbb465760521
	github.com/rs/zerolog v1.26.1
	github.com/shopspring/decimal v1.3.1
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	golang.org/x/crypto v0.0.0-20220321153916-2c7772ba3064
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
	sylr.dev/yaml/age/v3 v3.0.0-20211029141500-720306745b0b
	sylr.dev/yaml/v3 v3.0.0-20210127132132-941109e4f08c
)

replace github.com/quickfixgo/quickfix => github.com/sylr/quickfix v0.6.1-0.20220305100544-6fa3a91d941a

require (
	filippo.io/edwards25519 v1.0.0-rc.1 // indirect
	github.com/armon/go-proxyproto v0.0.0-20210323213023-7e956b284f0a // indirect
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/klauspost/compress v1.15.1 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/nats-io/jwt/v2 v2.2.1-0.20220113022732-58e87895b296 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/quickfixgo/fixt11 v0.0.0-20171007213433-d9788ca97f5d // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sys v0.0.0-20220319134239-a9b59b0215f8 // indirect
	golang.org/x/time v0.0.0-20220224211638-0e9765cccd65 // indirect
)
