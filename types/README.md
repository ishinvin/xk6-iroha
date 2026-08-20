# xk6-iroha-types

Ambient TypeScript declarations for [xk6-iroha](https://github.com/ishinvin/xk6-iroha)'s
`k6/x/iroha` module — gives editors IntelliSense on `signTransfer`,
`signGrantPermission`, and `reSignWithBatchMeta` in k6 scripts. Editor-only:
has no effect on how the extension is built or how scripts actually run.

## Install

```bash
pnpm add -D xk6-iroha-types
# or: npm install --save-dev xk6-iroha-types
```

## Usage

List it in your project's `jsconfig.json`/`tsconfig.json`:

```json
{
  "compilerOptions": {
    "types": ["xk6-iroha-types"]
  }
}
```

Or reference it directly in a script:

```js
/// <reference types="xk6-iroha-types" />
```

Hand-maintained, not generated — kept in sync with this repo's
[`transfers.go`](../transfers.go) and [`grants.go`](../grants.go) whenever
their signatures change. See [`xk6-iroha.d.ts`](./xk6-iroha.d.ts).
