package iroha

import (
	"encoding/base64"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/ishinvin/irohasign"
	irohapb "github.com/ishinvin/irohasign/proto"
)

// grantExport is signGrantPermission's return value: a base64-encoded,
// unsigned-batch-metadata transaction plus its reduced hash. Base64 is
// manual — a raw []byte returned to goja doesn't auto-marshal like
// encoding/json does.
type grantExport struct {
	Transaction string `js:"transaction"`
	ReducedHash string `js:"reducedHash"`
}

// SignGrantPermission builds and signs a grant-permission transaction from
// creatorAccountID to granteeAccountID, one GrantPermission command per
// entry in permissions (Iroha's GrantablePermission names, e.g.
// "can_add_my_signatory" — see irohasign/proto's GrantablePermission_value).
// The result carries no batch metadata; pass it to ReSignWithBatchMeta if
// it needs to join an atomic batch.
func (*Instance) SignGrantPermission(
	publicKey, privateKey, creatorAccountID, granteeAccountID string, permissions []string,
) (grantExport, error) {
	kp, err := irohasign.KeyPairFromHex(publicKey, privateKey)
	if err != nil {
		return grantExport{}, fmt.Errorf("iroha: parse key pair: %w", err)
	}

	b := irohasign.NewTransactionBuilder(creatorAccountID, uint64(time.Now().UnixMilli()))
	for _, p := range permissions {
		val, ok := irohapb.GrantablePermission_value[p]
		if !ok {
			return grantExport{}, fmt.Errorf("iroha: unknown GrantablePermission %q", p)
		}
		b = b.GrantPermission(granteeAccountID, irohapb.GrantablePermission(val))
	}

	reducedHash, err := b.ReducedHashHex()
	if err != nil {
		return grantExport{}, fmt.Errorf("iroha: reduced hash: %w", err)
	}
	tx, err := b.Sign(kp)
	if err != nil {
		return grantExport{}, fmt.Errorf("iroha: sign: %w", err)
	}
	txBytes, err := proto.Marshal(tx)
	if err != nil {
		return grantExport{}, fmt.Errorf("iroha: marshal: %w", err)
	}

	return grantExport{
		Transaction: base64.StdEncoding.EncodeToString(txBytes),
		ReducedHash: reducedHash,
	}, nil
}

// ReSignWithBatchMeta re-signs a base64-encoded transaction (e.g.
// SignGrantPermission's Transaction field) as part of an atomic batch,
// given the batch's ordered reduced hashes, and returns the result
// base64-encoded.
func (*Instance) ReSignWithBatchMeta(
	publicKey, privateKey, transactionBase64 string, reducedHashesHex []string,
) (string, error) {
	kp, err := irohasign.KeyPairFromHex(publicKey, privateKey)
	if err != nil {
		return "", fmt.Errorf("iroha: parse key pair: %w", err)
	}

	txBytes, err := base64.StdEncoding.DecodeString(transactionBase64)
	if err != nil {
		return "", fmt.Errorf("iroha: decode transaction: %w", err)
	}

	tx, err := irohasign.ReSignWithBatchMeta(txBytes, kp, reducedHashesHex)
	if err != nil {
		return "", fmt.Errorf("iroha: sign: %w", err)
	}

	signedBytes, err := proto.Marshal(tx)
	if err != nil {
		return "", fmt.Errorf("iroha: marshal: %w", err)
	}
	return base64.StdEncoding.EncodeToString(signedBytes), nil
}
