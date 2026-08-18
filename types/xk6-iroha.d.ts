// Ambient TypeScript declaration for the k6/x/iroha module this repo
// registers (see ../iroha.go's modules.Register call). Editor-only —
// gives IntelliSense to k6 scripts that import 'k6/x/iroha', with no
// effect on how the extension is built or how those scripts run.
//
// Hand-maintained, not generated: keep these signatures in sync by hand
// with ../transfers.go's SignTransfer and ../grants.go's
// SignGrantPermission / ReSignWithBatchMeta whenever those change. A Go
// (T, error) return with a non-nil error surfaces in JS as a thrown
// exception, not a second return value, so only the success shape is
// declared here.
declare module 'k6/x/iroha' {
  // transfers.go: SignTransfer — returns the signed TransferAsset
  // transaction, base64-encoded.
  export function signTransfer(
    srcPublicKey: string,
    srcPrivateKey: string,
    srcAccountId: string,
    dstAccountId: string,
    assetId: string,
    description: string,
    amount: string,
  ): string;

  // grants.go: SignGrantPermission — returns the signed, batch-metadata-
  // free grant transaction (base64) plus its reduced hash.
  export function signGrantPermission(
    publicKey: string,
    privateKey: string,
    creatorAccountId: string,
    granteeAccountId: string,
    permissions: string[],
  ): { transaction: string; reducedHash: string };

  // grants.go: ReSignWithBatchMeta — re-signs a base64-encoded
  // transaction (e.g. signGrantPermission's transaction) as part of an
  // atomic batch, given the batch's ordered reduced hashes; returns the
  // result base64-encoded.
  export function reSignWithBatchMeta(
    publicKey: string,
    privateKey: string,
    transactionBase64: string,
    reducedHashesHex: string[],
  ): string;
}
