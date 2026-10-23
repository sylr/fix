# fix

`fix` is a toolbox for the [Financial Information eXchange](https://www.onixs.biz/fix-protocol.html) protocol.

`fix` v2 is built on [fix2](https://github.com/sylr/fix2) instead of QuickFIX/Go. Its configuration
format changed: see [Configuration](#configuration). `fix` v1, built on QuickFIX/Go, is the `v0.x` tags.

## Bootstrap configuration

`fix` configuration allows you to define initiators, acceptors and sessions. Those
are wrapped in contexts that allow you to switch from one `initiator`/`acceptor` to another
with the `--context` option.

```shell
# Create a default configuration in $HOME/.fix/
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
(`store.type: sqlite3`); they keep them in memory otherwise (`store.type: memory`, the default).

If you don't provide `store.path` it defaults to `$HOME/.fix/initiator.db` and
`$HOME/.fix/acceptor.db`. Several sessions can share a database.

```shell
# Create databases for the different initiators and acceptors
fix init database
```

## Acceptor

The acceptor bundled in `fix` is a FIX5.0SP2 server that takes `NewOrderSingle`
messages, forwards them to an embedded NATS server and sends an `ExecutionReport`
message with an `OrdStatus` set to `0` (New). It is built with the `acceptor` build tag.

## Build from sources

`fix` requires a go toolchain >= 1.27 to be built from sources, and a C toolchain for the
SQLite store (`CGO_ENABLED=0` builds keep sessions in memory only).

```shell
# Build locally in the dist/ folder
make build
# Cross build for different os/arch in the dist/ folder
make crossbuild
# Build and install the fix binary into your $GOPATH/bin/ folder
make install
```

```shell
# Or install a released version with the go toolchain
go install sylr.dev/fix/v2@latest
```

### FIX dictionaries

`fix` validates messages with the FIXT 1.1 and FIX 5.0 SP2 dictionaries it is built with, and
sends and reads them with typed code [fix2gen](https://github.com/sylr/fix2) generates from them
(packages `fixt11` and `fix50sp2`). QuickFIX's FIX 5.0 SP2 dictionary is too large for fix2:
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
location by using the `--config` option. Values can be encrypted with [age](https://age-encryption.org)
for your SSH keys.

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
  # tls: {serverName: localhost, insecureSkipVerify: false, caFile: ca.pem}
  store:
    type: sqlite3
    # path: $HOME/.fix/initiator.db
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

Sessions also take:

| Key | Description |
| --- | --- |
| `senderLocationID`, `targetLocationID` | Location IDs of both sides. |
| `qualifier` | Tells apart sessions whose identities are the same on the wire. |
| `dictionaries.session`, `dictionaries.app` | QuickFIX XML dictionaries replacing the built-in ones. |
| `schedule` | When the session runs, starting afresh at sequence number 1 each window: daily from `start` to `end` (`HH:MM:SS`) on `days` (`mon`, `tuesday`…, every day if empty), or weekly from `startDay` at `start` to `endDay` at `end`; `location` is an IANA time zone, UTC by default. |
| `resetOnLogon`, `resetOnLogout`, `resetOnDisconnect` | Reset sequence numbers. |
| `logonTimeout`, `logoutTimeout` | Default to 5s, or `--timeout`. |
| `timestamps` | Precision of the timestamps sent: `seconds`, `milliseconds` (default), `microseconds`, `nanoseconds`. |
| `maxLatency` | How far a received SendingTime may be from our clock, 120s by default. |
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
| `ReconnectInterval` | `reconnectMin`, `reconnectMax` (initiators) |
| `HeartBtInt: 5` | `heartBtInt: 5s` |
| `TransportDataDictionary`, `AppDataDictionary` | optional `dictionaries: {session, app}` |
| `StartTime`, `EndTime`, `TimeZone` | `schedule: {start, end, location}`, daily |
| `StartDay`, `EndDay` | `schedule: {startDay, endDay}`, weekly |
| `RejectInvalidMessage: false` | `skipChecks: [...]` |
| other session keys | same name, lowerCamelCase (`SenderCompID` → `senderCompID`) |
| `--quickfix-logging` | `--fix-logging` |

`LogLimit` and the `postgres` store are gone.

v2 keeps its sessions in tables of its own: v1's sequence numbers, in the same `$HOME/.fix/*.db`
files, are not carried over, and a v2 session starts at sequence number 1. Unless your
counterparty resets its own, log on once with `resetOnLogon: true`, or agree on a reset.
