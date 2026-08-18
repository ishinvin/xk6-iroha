# xk6-iroha

[![CI](https://github.com/ishinvin/xk6-iroha/actions/workflows/ci.yml/badge.svg)](https://github.com/ishinvin/xk6-iroha/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A [k6](https://k6.io) extension that registers `k6/x/iroha`, a stateless
bridge onto [irohasign](https://github.com/ishinvin/irohasign) for signing
Hyperledger Iroha v1 transactions from load-test scripts.

## Why

k6's JS runtime (goja) can't build or sign Iroha protobuf transactions
itself — no Ed25519/protobuf support, and Iroha uses its own non-standard
Ed25519-SHA3 construction. This extension does that work in Go, in-process,
so a k6 script can get a signed, base64-encoded transaction back with a
single function call.

## What's here

- `iroha.go` — module registration and JS exports
- `transfers.go` — `signTransfer`
- `grants.go` — `signGrantPermission`, `reSignWithBatchMeta`

## Usage

```js
import iroha from "k6/x/iroha";

export default function () {
  const txBase64 = iroha.signTransfer(
    srcPublicKey,
    srcPrivateKey,
    srcAccountID,
    dstAccountID,
    assetID,
    description,
    amount,
  );

  const grant = iroha.signGrantPermission(
    publicKey,
    privateKey,
    creatorAccountID,
    granteeAccountID,
    ["can_add_my_signatory"],
  );
  // grant.transaction, grant.reducedHash

  const resigned = iroha.reSignWithBatchMeta(
    publicKey,
    privateKey,
    grant.transaction,
    [grant.reducedHash, otherReducedHash1, otherReducedHash2],
  );
}
```

## Editor types

`types/xk6-iroha.d.ts` declares `k6/x/iroha`'s three exports for
TypeScript-aware editors (VS Code's JS IntelliSense picks up ambient
`.d.ts` files automatically when this repo is visible in the same
window as the k6 scripts importing it — no build step or package.json
involved). It's hand-maintained, not generated, so update it by hand
alongside `transfers.go`/`grants.go` if their signatures change.

Release archives (see Install below) bundle this file alongside the
`k6-iroha` binary, so it's still available even without cloning this
repo — copy `types/xk6-iroha.d.ts` out of the downloaded archive into
wherever your editor is already looking (e.g. next to your k6 scripts)
to get the same IntelliSense.

## Install

Prebuilt `k6-iroha` binaries (a k6 build with this extension linked in) are
published on the
[Releases](https://github.com/ishinvin/xk6-iroha/releases) page for
linux/darwin (amd64+arm64) and windows/amd64.

Download and install one locally:

```
curl -LO https://github.com/ishinvin/xk6-iroha/releases/download/<tag>/xk6-iroha-<tag>-<os>-<arch>.tar.gz
tar -xzf xk6-iroha-<tag>-<os>-<arch>.tar.gz
sudo mv k6-iroha /usr/local/bin/
```

> **macOS:** the binary is quarantined on first download, so Gatekeeper will
> refuse to run it. Clear the flag once and you're set:
> `xattr -d com.apple.quarantine /usr/local/bin/k6-iroha`

Or build it yourself, either via [xk6](https://github.com/grafana/xk6):

```
xk6 build --with github.com/ishinvin/xk6-iroha@latest
```

or from a clone of this repo:

```
make build                          # -> bin/k6-iroha
sudo mv bin/k6-iroha /usr/local/bin/
```

Either way you get a `k6-iroha` binary — Run k6 scripts with it exactly as you would with `k6`: `k6-iroha run script.js`.

## Development

```
make test    # run tests
make lint    # run golangci-lint
make build   # build a k6 binary with this extension, via xk6
make setup   # install tools and git hooks
```
