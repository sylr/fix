# fix

`fix` is a toolbox for the [Financial Information eXchange](https://www.onixs.biz/fix-protocol.html) protocol.

## Bootsrap configuration

`fix` configuration allows you to define initiators, acceptors and sessions. Those
are wrapped in contexts that allow you to switch from one `initiator`/`acceptor` to another
with the `--context` option.

```shell
# Create a default configuration in $HOME/.fix/
fix config init
```

An `initiator` describes the network properties needed to connect to a FIX acceptor.
An `acceptor` describes the network properites needed to launch a server which will
accept connection from FIX initiators.
A `session` describes the properties of a FIX session.

An `acceptor` can have several sessions but an `initiator` can only have one.

You can use the `current-context` property to define the context to use if none is
provided as an option.

## Boostrap databases

The FIX protocol requires a message sequence number which should be incremented for
each messages sent on the wire. In order to keep track of the sequence number `fix`
uses databases. Note that you can share a database for different initiators and another
for the acceptors but you can't share a database between initiators and acceptors.

**Currently only the `sqlite3` driver is supported** but we plan to support for 
other database types such as `mysql` & `postgresql`.

You can select the driver using `SQLStoreDriver` and you can give a connection string
using the `SQLStoreDataSourceName` configuration in the initiators and acceptors.

If you don't provide `SQLStoreDataSourceName` for the `sqlite3` driver it will default
to `$HOME/.fix/initiator.db` & `$HOME/.fix/acceptor.db`. If you want to use a different
path or custom options refer to the [following documentation](https://github.com/mattn/go-sqlite3#usage-1).

```shell
# Create databases for the different initiators and acceptors
fix database init
```

## Acceptor

The acceptor bundled in `fix` is a FIX5.0SP2 server that takes `NewSingleOrder`
messages, forward them to a embeded NATS server and send an `ExecutionReportStatus`
message with and `OrdStatus` set to `0` (New).

## Build from sources

`fix` requires a go toolchain >= 1.18 to be built from sources. You'll also require `libsqlite3`.

```shell
# Build locally in the dist/ folder
make build
# Cross build for different os/arch in the dist/ folder
make crossbuild
# Build and install the fix binary into your $GOPATH/bin/ folder
make install
```

## Configuration

Default configuration is located at `$HOME/.fix/config`. You can specify a custom
location by using the `--config` option.

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
  SocketAcceptHost: 127.0.0.1
  SocketAcceptPort: 5005
  SQLStoreDriver: sqlite3
initiators:
- name: localhost
  SocketConnectHost: 127.0.0.1
  SocketConnectPort: 5005
  SocketServerName: localhost
  SocketUseSSL: false
  SocketInsecureSkipVerify: false
  SocketTimeout: 5s
  SQLStoreDriver: sqlite3
sessions:
- name: server
  HeartBtInt: 5
  SenderCompID: BIGCORP
  TargetCompID: smallcorp
  TargetSubID: john
  BeginString: FIXT.1.1
  DefaultApplVerID: FIX.5.0SP2
  TransportDataDictionary: $HOME/.fix/FIXT11.xml
  AppDataDictionary: $HOME/.fix/FIXT11.xml
- name: localhost
  HeartBtInt: 5
  SenderCompID: smallcorp
  SenderSubID: john
  TargetCompID: BIGCORP
  BeginString: FIXT.1.1
  DefaultApplVerID: FIX.5.0SP2
  TransportDataDictionary: $HOME/.fix/FIXT11.xml
  AppDataDictionary: $HOME/.fix/FIX50SP2.xml
```
