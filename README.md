# fix

`fix` is a toolbox for the [Financial Information eXchange](https://www.onixs.biz/fix-protocol.html) protocol.

`fix` v2 is built on [fix2](https://github.com/sylr/fix2) instead of QuickFIX/Go. Its configuration
format changed: see [Configuration](#configuration). `fix` v1, built on QuickFIX/Go, is the `v0.x` tags.

## FIX versions

`fix` speaks FIX 5.0 SP2 over FIXT 1.1 (`beginString: FIXT.1.1`, `defaultApplVerID: FIX.5.0SP2`),
FIX 4.4 (`beginString: FIX.4.4`) and FIX 4.3 (`beginString: FIX.4.3`). Commands refuse other versions,
and refuse to send a message the session's version does not define, or that lacks a field it
requires, rather than have the counterparty reject it: options such as `--origination` or
`--attribute-type`, whose fields came with FIX 5.0, only work with FIX 5.0 SP2 sessions, and
`--party-sub-ids`, whose typed sub IDs came with FIX 4.4, only with FIX 4.4 and later. With
FIX 4.3, orders carry `HandlInst=1` (automated execution, no broker intervention), which it requires.

## Bootstrap configuration

`fix` configuration allows you to define initiators, acceptors and sessions. Those
are wrapped in contexts that allow you to switch from one `initiator`/`acceptor` to another
with the `--context` option.

```shell
# Create a default configuration in $HOME/.fix/config, unless it exists
fix init config
```

An `initiator` describes the network properties needed to connect to a FIX acceptor.
An `acceptor` describes the network properties needed to launch a server which will
accept connection from FIX initiators.
A `session` describes the properties of a FIX session.

⚠️ An `acceptor` can have several sessions but an <ins>`initiator` can only have one</ins>.

You can use the `current-context` property to define the context to use if none is
provided as an option.

## Bootstrap databases

The FIX protocol requires a message sequence number which should be incremented for
each message sent on the wire. In order to keep track of the sequence numbers across
runs, initiators and acceptors can keep their sessions in a SQLite database
(`store.type: sqlite3`, which `fix init config` writes unless `fix` is built without cgo).
They keep them in memory otherwise (`store.type: memory`, the default): each run then starts
at sequence number 1, which counterparties refuse unless sequence numbers are reset, as
`resetOnLogon: true` does.

If you don't provide `store.path` it defaults to `$HOME/.fix/initiator.db` and
`$HOME/.fix/acceptor.db`. Several initiators and acceptors can share a database.

A database and its tables are created when first used, so this step is optional:

```shell
# Create the databases of the initiators and acceptors, with the sessions of their contexts.
# --overwrite replaces existing ones, starting their sessions afresh at sequence number 1.
fix init database
```

## Acceptor

The acceptor bundled in `fix` is a FIX 5.0 SP2, 4.4 and 4.3 server that takes `NewOrderSingle`
messages, forwards them to an embedded NATS server and sends an `ExecutionReport`
message with an `OrdStatus` set to `0` (New). It is built with the `acceptor` build tag.

## Build from sources

`fix` requires a go toolchain >= 1.27 to be built from sources, and a C toolchain for the
SQLite store (`CGO_ENABLED=0` builds keep sessions in memory only). The `acceptor` and `validator`
commands are only built with the `acceptor` and `validator` build tags, which `make` enables by
default (`GO_BUILD_TAGS`).

```shell
# Build locally in the dist/ folder
make build
# Cross build for different os/arch in the dist/ folder
make crossbuild
# Build and install the fix binary into your $GOPATH/bin/ folder
make install
```

```shell
# Or install a tagged v2 release with the go toolchain
go install -tags acceptor,validator sylr.dev/fix/v2@latest
```

### FIX dictionaries

`fix` validates messages with the FIXT 1.1, FIX 5.0 SP2, FIX 4.4 and FIX 4.3 dictionaries it is
built with, and sends and reads them with typed code [fix2gen](https://github.com/sylr/fix2)
generates from the FIXT 1.1 and FIX 5.0 SP2 ones (packages `fixt11` and `fix50sp2`), whose tags
and values FIX 4.3 and 4.4 share. The FIX 4.3 and 4.4 dictionaries are QuickFIX's, whole
(packages `fix43` and `fix44`). QuickFIX's FIX 5.0 SP2 dictionary is too large for fix2:
`fix50sp2/FIX50SP2.xml` keeps only the messages `fix` uses, without the components describing
legs, underlyings and derivatives terms. Sessions using it accept the tags of those components
(`skipChecks: [unknownTags]` is implied).

```shell
# Regenerate the dictionaries and typed code, e.g. after adding a MsgType in fix50sp2/generate.go
make generate
```

A session can use other QuickFIX XML dictionaries with `dictionaries.session` and `dictionaries.app`.

## Configuration

Default configuration is located at `$HOME/.fix/config`. You can specify a custom
location by using the `--config` option. Unknown keys are refused. Values can be encrypted
with [age](https://age-encryption.org) for your SSH keys (`~/.ssh/id_rsa*`, `~/.ssh/id_ed25519*`)
and tagged `!crypto/age`; keys protected by a passphrase are only used with `--interactive`,
the default. Paths (`store.path`, `dictionaries`, `tls` files) may hold environment variables
and start with `~/`.

```yaml
# vim: syntax=yaml :
---
current-context: localhost
contexts:
- name: localhost
  initiator: localhost
  sessions: [localhost]
- name: server
  acceptor: server
  sessions: [server]
acceptors:
- name: server
  address: 127.0.0.1:5005
  # tls: {certFile: server.pem, keyFile: server.key, caFile: ca.pem, clientAuth: true}
  store:
    type: sqlite3
initiators:
- name: localhost
  addresses: [127.0.0.1:5005]  # tried in turn
  dialTimeout: 5s
  # tls: {serverName: localhost, caFile: ca.pem}
  store:
    type: sqlite3
    # path: ~/.fix/initiator.db
sessions:
- name: localhost
  beginString: FIXT.1.1
  defaultApplVerID: FIX.5.0SP2
  senderCompID: smallcorp
  senderSubID: john
  targetCompID: BIGCORP
  heartBtInt: 5s
  # username: john
  # password: secret
- name: server
  beginString: FIXT.1.1
  defaultApplVerID: FIX.5.0SP2
  senderCompID: BIGCORP
  targetCompID: smallcorp
  targetSubID: john
  heartBtInt: 5s
```

Contexts take a `name`, an `initiator` or an `acceptor`, and the names of their `sessions`.

Initiators take:

| Key | Description |
| --- | --- |
| `name` | Required. |
| `addresses` | Required: the acceptor's `host:port` addresses, tried in turn. |
| `dialTimeout` | Bounds each connection attempt: `--timeout` if given, else this, else 5s. |
| `reconnectMin`, `reconnectMax` | Bound the wait between connection attempts, 1s and 30s by default; `reconnectMin` may not be above `reconnectMax`. |
| `tls` | Makes connections TLS: `serverName` the acceptor's certificate is checked against (the host dialed by default), `insecureSkipVerify` to accept any, `caFile` holding the authorities trusted (the system's by default), and `certFile` and `keyFile`, a client certificate. |
| `store` | `type`: `memory` (default) or `sqlite3`, whose database is `path`, `$HOME/.fix/initiator.db` by default. |

Acceptors take:

| Key | Description |
| --- | --- |
| `name` | Required. |
| `address` | Required: the `host:port` address to listen on, `:port` for every interface. |
| `tls` | Makes connections TLS: `certFile` and `keyFile` are required, and `clientAuth: true` requires client certificates, verified against the authorities in `caFile`, which it requires. |
| `store` | As for initiators, `$HOME/.fix/acceptor.db` by default. |

Durations are written `5s`, `1m30s`… and may not be negative, but for `maxLatency`.

Sessions take:

| Key | Description |
| --- | --- |
| `name`, `senderCompID`, `targetCompID` | Required. |
| `beginString` | Required: `FIXT.1.1`, `FIX.4.4` or `FIX.4.3`. |
| `defaultApplVerID` | Required for `FIXT.1.1`, and only for it: `FIX.5.0SP2`, or its ApplVerID `9`. |
| `senderSubID`, `senderLocationID`, `targetSubID`, `targetLocationID` | Sub and location IDs of both sides. |
| `qualifier` | Never sent: tells apart, in a store, sessions whose identities on the wire are the same, so that they keep their own sequence numbers. An acceptor refuses two sessions only their qualifier tells apart. |
| `heartBtInt` | Heartbeat interval an initiator asks for, whole seconds, 30s by default; an acceptor uses the initiator's. |
| `username`, `password` | Username (553) and Password (554) of the Logons sent. |
| `dictionaries.session`, `dictionaries.app` | QuickFIX XML dictionaries replacing the built-in ones. A FIX 4.x session has one dictionary: set either, not both. |
| `schedule` | When the session runs, starting afresh at sequence number 1 each window: daily from `start` to `end` (`HH:MM:SS`), every day or on the `days` (`mon`, `tuesday`…) windows start on, or weekly from `startDay` at `start` to `endDay` at `end`. An `end` before or equal to `start` ends the next day, or the next week for a weekly schedule starting and ending the same day: `start` equal to `end` runs a whole day, or week. `location` is an IANA time zone, UTC by default. |
| `resetOnLogon`, `resetOnLogout`, `resetOnDisconnect` | Reset sequence numbers. |
| `logonTimeout`, `logoutTimeout` | `--timeout` if given, else these, else 5s. |
| `timestamps` | Precision of the timestamps sent: `seconds`, `milliseconds` (default), `microseconds`, `nanoseconds`. |
| `maxLatency` | How far a received SendingTime may be from our clock, 120s by default; negative disables the check. |
| `nextExpectedMsgSeqNum`, `lastMsgSeqNumProcessed` | Send and honour NextExpectedMsgSeqNum (789), send LastMsgSeqNumProcessed (369). |
| `tolerations` | Departures from the FIX specifications to allow: `logonWithoutEncryptMethod`, `unsupportedDefaultApplVerID`, `logoutBeforeGapRecovery`, `sequenceResetWithoutOrigSendingTime`, `resetLogonSeqNumNotOne`, `malformedNextExpectedMsgSeqNum`, `noHeartbeatAwaitingTestRequest`, `tagLeadingZeros`, `finerTimestamps`, `emptyRequiredGroups`, or `quickfix` for those QuickFIX engines need. |
| `skipChecks` | Dictionary validations to skip: `emptyValues`, `formats`, `repeatedTags`, `order`, `groupCounts`, `unknownTags`, `userDefinedTags`, `groupRequired`, `groupOrder`, `timestampPrecision`, `emptyRequiredGroups`. |

### Migrating from v1

| v1 | v2 |
| --- | --- |
| `SocketConnectHost`, `SocketConnectPort` | `addresses: [host:port]` |
| `SocketAcceptHost`, `SocketAcceptPort` | `address: host:port` |
| `SocketUseSSL`, `SocketServerName`, `SocketInsecureSkipVerify`, `SocketCAFile`, `SocketCertificateFile`, `SocketPrivateKeyFile` | `tls: {serverName, insecureSkipVerify, caFile, certFile, keyFile}` |
| `SocketTimeout` | `dialTimeout` (initiators), `logonTimeout` / `logoutTimeout` (sessions) |
| `SQLStoreDriver`, `SQLStoreDataSourceName` | `store: {type: sqlite3, path}` |
| `ReconnectInterval` (sessions) | `reconnectMin`, `reconnectMax` (initiators) |
| `HeartBtInt: 5` | `heartBtInt: 5s` |
| `TransportDataDictionary`, `AppDataDictionary` | optional `dictionaries: {session, app}` |
| `StartTime`, `EndTime`, `TimeZone` | `schedule: {start, end, location}`, daily |
| `StartDay`, `EndDay` | `schedule: {startDay, endDay}`, weekly |
| `RejectInvalidMessage: false` | `skipChecks: [...]` (sessions) |
| `SessionQualifier` | `qualifier` |
| other session keys | same name, lowerCamelCase (`SenderCompID` → `senderCompID`) |
| `--quickfix-logging` | `--fix-logging` |

`LogLimit`, `UseTCPProxy` and the `postgres` store are gone.

Commands changed too:

| v1 | v2 |
| --- | --- |
| `status order --origclordid` | `--clordid`, sent as ClOrdID (11), the order's own; OrigClOrdID (41) is never sent |
| `--copy-credentials-fom-config` | `--copy-credentials-from-config` |
| `status tradingsession --trading-session-id`, sent as TradSesReqID (335) | `--request-id`; `--trading-session-id` now sets TradingSessionID (336) |
| `list security`, SecurityReqID always `4` | a generated SecurityReqID or `--request-id`; `--symbol` for `--type symbol`, whose default is now `all_securities` |
| `new quote --origination`, `--attribute-type`, `--attribute-value` | gone: a Quote carries neither in any FIX version |
| `acceptor --no-nats-embeded` | `acceptor --nats-embedded=false` |
| `--price`, `--quantity` rounded to 2 decimals | sent exactly as given |
| `--initiator`/`--session` ignored when `current-context` is set | they take precedence over it |

Commands exit with a non-zero status when the counterparty rejects their request (Reject or
BusinessMessageReject), logs out, or does not answer in time, and log to standard error, standard
output holding what they received only.

v2 keeps its sessions in tables of its own: v1's sequence numbers, in the same `$HOME/.fix/*.db`
files, are not carried over, and a v2 session starts at sequence number 1. Unless your
counterparty resets its own, log on once with `resetOnLogon: true`, or agree on a reset.
