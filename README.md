# fix

## Install

```
gotip install -trimpath sylr.dev/fix@latest
```

## Configuration

Configuration must be located at `$HOME/.fix`.

```yaml
---
contexts:
- name: dev
  acceptor: eks-euc1-dev-02
  session: sylvain
acceptors:
- name: company-cleartext
  SocketConnectHost: fix.elb.eu-central-1.amazonaws.com
  SocketConnectPort: 5005
  SocketUseSSL: false
  SocketInsecureSkipVerify: false
  SocketTimeout: 5s
- name: company-tls
  SocketConnectHost: fix-tls.elb.eu-central-1.amazonaws.com
  SocketConnectPort: 5005
  SocketServerName: fix.company.com
  SocketUseSSL: true
  SocketInsecureSkipVerify: false
  SocketTimeout: 5s
sessions:
- name: sylvain
  HeartBtInt: 5
  TargetCompID: COMPANY
  SenderCompID: CLIENT
  BeginString: FIXT.1.1
  DefaultApplVerID: FIX.5.0SP2
  Username: xxxxxxxxxx
```
