module sylr.dev/fix

go 1.22

require (
	filippo.io/age v1.2.1
	github.com/dustin/go-humanize v1.0.1
	github.com/google/uuid v1.6.0
	github.com/iancoleman/strcase v0.3.0
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/nats-io/nats-server/v2 v2.10.24
	github.com/nats-io/nats.go v1.38.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/prometheus/client_golang v1.20.5
	github.com/quickfixgo/enum v0.1.0
	github.com/quickfixgo/field v0.1.0
	github.com/quickfixgo/fixt11 v0.1.0
	github.com/quickfixgo/quickfix v0.9.6
	github.com/quickfixgo/tag v0.1.0
	github.com/rs/zerolog v1.33.0
	github.com/shopspring/decimal v1.4.0
	github.com/spf13/cobra v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/sylr/quickfixgo-fix50sp2/marketdataincrementalrefresh v0.0.0-20220401195242-281940b8a21e
	github.com/sylr/quickfixgo-fix50sp2/marketdatasnapshotfullrefresh v0.0.0-20220401195242-281940b8a21e
	golang.org/x/crypto v0.32.0
	golang.org/x/term v0.28.0
	sylr.dev/yaml/age/v3 v3.0.0-20221203153010-eb6b46db8d90
	sylr.dev/yaml/v3 v3.0.0-20220527135632-500fddf2b049
)

replace (
	github.com/quickfixgo/enum => github.com/sylr/quickfixgo-enum v0.0.0-20220401193143-29a559514373
	github.com/quickfixgo/field => github.com/sylr/quickfixgo-field v0.0.0-20220401193046-ca4cd16301d2
	github.com/quickfixgo/quickfix => github.com/sylr/quickfix-go v0.9.7-0.20250122164305-17249b610276
	github.com/quickfixgo/tag => github.com/sylr/quickfixgo-tag v0.0.0-20220401193001-96cf7367fdfa
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/minio/highwayhash v1.0.3 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nats-io/jwt/v2 v2.7.3 // indirect
	github.com/nats-io/nkeys v0.4.9 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/pires/go-proxyproto v0.8.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.62.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/smartystreets/assertions v1.13.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	google.golang.org/protobuf v1.36.3 // indirect
)
