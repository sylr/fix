module sylr.dev/fix

go 1.18

require (
	filippo.io/age v1.1.1
	github.com/dustin/go-humanize v1.0.0
	github.com/google/uuid v1.3.0
	github.com/iancoleman/strcase v0.2.0
	github.com/lib/pq v1.10.7
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/nats-io/nats-server/v2 v2.9.10
	github.com/nats-io/nats.go v1.22.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/prometheus/client_golang v1.14.0
	github.com/quickfixgo/enum v0.0.0-20210629025633-9afc8539baba
	github.com/quickfixgo/field v0.0.0-20171007195410-74cea5ec78c7
	github.com/quickfixgo/fixt11 v0.0.0-20171007213433-d9788ca97f5d
	github.com/quickfixgo/quickfix v0.6.1-0.20210618140103-31f5ebe90229
	github.com/quickfixgo/tag v0.0.0-20171007194743-cbb465760521
	github.com/rs/zerolog v1.28.0
	github.com/shopspring/decimal v1.3.1
	github.com/spf13/cobra v1.6.1
	github.com/spf13/pflag v1.0.5
	github.com/sylr/quickfixgo-fix50sp2/marketdataincrementalrefresh v0.0.0-20220401195242-281940b8a21e
	github.com/sylr/quickfixgo-fix50sp2/marketdatasnapshotfullrefresh v0.0.0-20220401195242-281940b8a21e
	golang.org/x/crypto v0.4.0
	golang.org/x/term v0.3.0
	sylr.dev/yaml/age/v3 v3.0.0-20221203153010-eb6b46db8d90
	sylr.dev/yaml/v3 v3.0.0-20220527135632-500fddf2b049
)

replace (
	github.com/quickfixgo/enum => github.com/sylr/quickfixgo-enum v0.0.0-20220401193143-29a559514373
	github.com/quickfixgo/field => github.com/sylr/quickfixgo-field v0.0.0-20220401193046-ca4cd16301d2
	github.com/quickfixgo/quickfix => github.com/sylr/quickfix-go v0.6.1-0.20221223080000-9e4d31ed9df6
	github.com/quickfixgo/tag => github.com/sylr/quickfixgo-tag v0.0.0-20220401193001-96cf7367fdfa
)

require (
	filippo.io/edwards25519 v1.0.0 // indirect
	github.com/armon/go-proxyproto v0.0.0-20210323213023-7e956b284f0a // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.15.13 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/minio/highwayhash v1.0.2 // indirect
	github.com/nats-io/jwt/v2 v2.3.0 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.39.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/rivo/uniseg v0.4.3 // indirect
	github.com/rogpeppe/go-internal v1.8.1 // indirect
	github.com/smartystreets/assertions v1.13.0 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
