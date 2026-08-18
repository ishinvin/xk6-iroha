package iroha

import (
	"encoding/base64"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/ishinvin/irohasign"
	irohapb "github.com/ishinvin/irohasign/proto"
)

func TestSignTransfer(t *testing.T) {
	kp, err := irohasign.GenerateKeyPair()
	if err != nil {
		t.Fatalf("GenerateKeyPair: %v", err)
	}
	inst := &Instance{}

	encoded, err := inst.SignTransfer(kp.PublicHex(), kp.PrivateHex(), "alice@test", "bob@test", "khr#test", "test transfer", "10.00")
	if err != nil {
		t.Fatalf("SignTransfer: %v", err)
	}

	txBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("result is not valid base64: %v", err)
	}
	var tx irohapb.Transaction
	if err := proto.Unmarshal(txBytes, &tx); err != nil {
		t.Fatalf("result is not a valid Transaction: %v", err)
	}

	if len(tx.GetSignatures()) != 1 {
		t.Fatalf("signatures = %d, want 1", len(tx.GetSignatures()))
	}
	commands := tx.GetPayload().GetReducedPayload().GetCommands()
	if len(commands) != 1 {
		t.Fatalf("commands = %d, want 1", len(commands))
	}
	transfer := commands[0].GetTransferAsset()
	if transfer == nil {
		t.Fatalf("command is not a TransferAsset")
	}
	if transfer.GetSrcAccountId() != "alice@test" || transfer.GetDestAccountId() != "bob@test" ||
		transfer.GetAssetId() != "khr#test" || transfer.GetAmount() != "10.00" {
		t.Fatalf("TransferAsset fields mismatch: %+v", transfer)
	}
}

func TestSignTransfer_BadKey(t *testing.T) {
	inst := &Instance{}
	_, err := inst.SignTransfer("not-hex", "not-hex", "alice@test", "bob@test", "khr#test", "", "1")
	if err == nil {
		t.Fatal("expected error for invalid key hex, got nil")
	}
	if !strings.Contains(err.Error(), "parse key pair") {
		t.Fatalf("error = %q, want it to mention key parsing", err)
	}
}
