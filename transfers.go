package iroha

import (
	"encoding/base64"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/ishinvin/irohasign"
)

// SignTransfer builds and signs a single TransferAsset transaction from
// srcAccountID to dstAccountID, using the given keypair, and returns the
// result base64-encoded. Stateless — the caller supplies the keypair
// directly, no prior "load" call needed.
func (*Instance) SignTransfer(
	srcPublicKey, srcPrivateKey, srcAccountID, dstAccountID, assetID, description, amount string,
) (string, error) {
	kp, err := irohasign.KeyPairFromHex(srcPublicKey, srcPrivateKey)
	if err != nil {
		return "", fmt.Errorf("iroha: parse key pair: %w", err)
	}

	tx, err := irohasign.NewTransactionBuilder(srcAccountID, uint64(time.Now().UnixMilli())).
		TransferAsset(srcAccountID, dstAccountID, assetID, description, amount).
		Sign(kp)
	if err != nil {
		return "", fmt.Errorf("iroha: sign transfer %s -> %s: %w", srcAccountID, dstAccountID, err)
	}
	txBytes, err := proto.Marshal(tx)
	if err != nil {
		return "", fmt.Errorf("iroha: marshal transfer %s -> %s: %w", srcAccountID, dstAccountID, err)
	}

	return base64.StdEncoding.EncodeToString(txBytes), nil
}
